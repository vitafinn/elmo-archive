package util

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func readLocalPublicCertificate() (*x509.Certificate, error) {
	data, err := os.ReadFile("ca.crt")
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("cannot read local certificate")
	}
	return x509.ParseCertificate(block.Bytes)
}

// Read-only status: unlike preparation this never generates a certificate.
func CaptureCertificateStatus() (*CaptureCertificate, error) {
	certificateMutex.Lock()
	defer certificateMutex.Unlock()
	cert, err := readLocalPublicCertificate()
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	result := certificateDetails(cert)
	store, err := openCertStore()
	if err != nil {
		return nil, err
	}
	defer closeCertStore(store)
	found, err := findExactCertificate(store, cert)
	if err != nil {
		return nil, err
	}
	if found != nil {
		result.Trusted = true
		windows.CertFreeCertificateContext(found)
	} else {
		result.Trusted = IsTrustedCA()
	}
	return &result, nil
}

func findExactCertificate(store windows.Handle, cert *x509.Certificate) (*windows.CertContext, error) {
	context, err := windows.CertCreateCertificateContext(windows.X509_ASN_ENCODING, &cert.Raw[0], uint32(len(cert.Raw)))
	if err != nil {
		return nil, err
	}
	defer windows.CertFreeCertificateContext(context)
	found, err := windows.CertFindCertificateInStore(store, windows.X509_ASN_ENCODING, 0, windows.CERT_FIND_EXISTING, unsafe.Pointer(context), nil)
	if errors.Is(err, syscall.Errno(windows.CRYPT_E_NOT_FOUND)) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(unsafe.Slice(found.EncodedCert, int(found.Length)), cert.Raw) {
		windows.CertFreeCertificateContext(found)
		return nil, errors.New("stored certificate differs; no trust was removed")
	}
	return found, nil
}

func removeExactCertificate(store windows.Handle, cert *x509.Certificate) error {
	found, err := findExactCertificate(store, cert)
	if err != nil {
		return err
	}
	if found == nil {
		return nil
	}
	// This API releases the returned context on both success and failure.
	return windows.CertDeleteCertificateFromStore(found)
}

func RemoveReviewedCaptureCertificate(thumbprint string) error {
	certificateMutex.Lock()
	defer certificateMutex.Unlock()
	cert, err := readLocalPublicCertificate()
	if err != nil {
		return err
	}
	if !cert.IsCA || (cert.Subject.CommonName != "Elmo Archive Local Capture CA" && cert.Subject.CommonName != "ElmoBeacon Custom CA") || certificateDetails(cert).Thumbprint != thumbprint {
		return errors.New("certificate identity changed or is not this app's certificate; review it again")
	}
	store, err := openCertStore()
	if err != nil {
		return err
	}
	defer closeCertStore(store)
	finishPrompt := beginCertificatePrompt()
	defer finishPrompt()
	if err = removeExactCertificate(store, cert); err != nil {
		return err
	}
	if IsTrustedCA() {
		return errors.New("current-user trust was removed, but Windows still trusts this certificate through another store; check Windows certificate settings")
	}
	return nil
}
