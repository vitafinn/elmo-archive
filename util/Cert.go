package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
)

var (
	crypt32                              = syscall.NewLazyDLL("crypt32.dll")
	procCertOpenStore                    = crypt32.NewProc("CertOpenStore")
	procCertCloseStore                   = crypt32.NewProc("CertCloseStore")
	procCertAddEncodedCertificateToStore = crypt32.NewProc("CertAddEncodedCertificateToStore")
	certificateMutex                     sync.Mutex
	certificatePromptActive              atomic.Bool
)

func GenCA() error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "Elmo Archive Local Capture CA",
			Organization: []string{"Elmo Archive"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(20, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}
	certOut, err := os.Create("ca.crt")
	if err != nil {
		return err
	}
	defer certOut.Close()
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return err
	}

	keyOut, err := os.Create("ca.key")
	if err != nil {
		return err
	}
	defer keyOut.Close()
	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		return err
	}

	return nil
}

func IsTrustedCA() bool {
	certPem, err := os.ReadFile("ca.crt")
	if err != nil {
		fmt.Println("Failed to read ca.crt:", err)
		return false
	}
	certDer, _ := pem.Decode(certPem)
	if certDer == nil {
		fmt.Println("PEM decode failed")
		return false
	}
	cert, err := x509.ParseCertificate(certDer.Bytes)
	if err != nil {
		fmt.Println("Failed to parse cert:", err)
		return false
	}

	_, err = cert.Verify(x509.VerifyOptions{})
	if err != nil {
		fmt.Println("Certificate verification failed:", err)
		return false
	}

	return true
}

func openCertStore() (windows.Handle, error) {
	storeName, _ := syscall.UTF16PtrFromString("Root")

	r, _, err := procCertOpenStore.Call(uintptr(10), 0, 0, uintptr(0x00010000), uintptr(unsafe.Pointer(storeName)))
	if r == 0 {
		return 0, fmt.Errorf("CertOpenStore error: %v", err)
	}
	return windows.Handle(r), nil
}

func closeCertStore(store windows.Handle) {
	if store != 0 {
		procCertCloseStore.Call(uintptr(store), 0)
	}
}

func installCertificateDER(der []byte) error {
	store, err := openCertStore()
	if err != nil {
		return err
	}
	defer closeCertStore(store)

	r, _, err := procCertAddEncodedCertificateToStore.Call(
		uintptr(store),
		uintptr(windows.X509_ASN_ENCODING),
		uintptr(unsafe.Pointer(&der[0])),
		uintptr(len(der)),
		uintptr(4),
		0,
	)

	if r == 0 {
		return fmt.Errorf("CertAddEncodedCertificateToStore error: %v", err)
	}
	return nil
}

// Only public certificate information is exposed to the UI, never the key.
type CaptureCertificate struct {
	Directory  string
	Name       string
	Thumbprint string
	Trusted    bool
	Legacy     bool
}

func certificateDetails(cert *x509.Certificate) CaptureCertificate {
	// SHA-1 is used only to match Windows' displayed identifier, not for signing.
	sum := sha1.Sum(cert.Raw)
	hex := fmt.Sprintf("%X", sum)
	groups := make([]string, 0, 5)
	for i := 0; i < len(hex); i += 8 {
		groups = append(groups, hex[i:i+8])
	}
	directory, _ := filepath.Abs(".")
	return CaptureCertificate{Directory: directory, Name: cert.Subject.CommonName, Thumbprint: strings.Join(groups, " "), Legacy: cert.Subject.CommonName == "ElmoBeacon Custom CA"}
}

func ensureLocalCertificate() (*tls.Certificate, error) {
	_, certErr := os.Stat("ca.crt")
	_, keyErr := os.Stat("ca.key")
	if os.IsNotExist(certErr) && os.IsNotExist(keyErr) {
		if err := GenCA(); err != nil {
			return nil, err
		}
	}
	pair, err := ParseCA()
	if err != nil {
		return nil, errors.New("local certificate files are missing or do not match; restore ca.crt and ca.key together, or use a new empty app folder")
	}
	cert := pair.Leaf
	if !cert.IsCA || time.Now().Before(cert.NotBefore) || time.Now().After(cert.NotAfter) {
		return nil, errors.New("the local capture certificate is invalid or expired; use a new empty app folder to create a replacement")
	}
	return pair, nil
}

func PrepareCaptureCertificate() (CaptureCertificate, error) {
	certificateMutex.Lock()
	defer certificateMutex.Unlock()
	pair, err := ensureLocalCertificate()
	if err != nil {
		return CaptureCertificate{}, err
	}
	result := certificateDetails(pair.Leaf)
	result.Trusted = IsTrustedCA()
	return result, nil
}

func InstallReviewedCertificate(thumbprint string) error {
	certificateMutex.Lock()
	defer certificateMutex.Unlock()
	pair, err := ParseCA()
	if err != nil {
		return errors.New("certificate changed; review certificate setup again")
	}
	if certificateDetails(pair.Leaf).Thumbprint != thumbprint {
		return errors.New("certificate changed; review its new thumbprint before installing")
	}
	finishPrompt := beginCertificatePrompt()
	defer finishPrompt()
	// Install the exact DER that was reviewed, even if the disk file changes.
	if err := installCertificateDER(pair.Leaf.Raw); err != nil {
		return errors.New("Windows did not install the certificate; you can cancel setup or try again")
	}
	if !IsTrustedCA() {
		return errors.New("Windows has not trusted the local certificate; capture has not started")
	}
	return nil
}

func ParseCA() (*tls.Certificate, error) {
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		return nil, err
	}

	caKey, err := os.ReadFile("ca.key")
	if err != nil {
		return nil, err
	}

	parsedCert, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return nil, err
	}
	if parsedCert.Leaf, err = x509.ParseCertificate(parsedCert.Certificate[0]); err != nil {
		return nil, err
	}
	return &parsedCert, nil
}
