package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"golang.org/x/sys/windows"
	"math/big"
	"os"
	"testing"
	"time"
)

func TestRemovalTouchesOnlyExactCertificateInMemory(t *testing.T) {
	store, err := windows.CertOpenStore(windows.CERT_STORE_PROV_MEMORY, 0, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer windows.CertCloseStore(store, 0)
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	certs := make([]*x509.Certificate, 2)
	for i := range certs {
		template := &x509.Certificate{SerialNumber: big.NewInt(int64(i + 1)), Subject: pkix.Name{CommonName: "Elmo Archive Local Capture CA"}, NotBefore: time.Now(), NotAfter: time.Now().Add(time.Hour), IsCA: true, BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign}
		der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
		if err != nil {
			t.Fatal(err)
		}
		certs[i], err = x509.ParseCertificate(der)
		if err != nil {
			t.Fatal(err)
		}
		context, err := windows.CertCreateCertificateContext(windows.X509_ASN_ENCODING, &der[0], uint32(len(der)))
		if err != nil {
			t.Fatal(err)
		}
		err = windows.CertAddCertificateContextToStore(store, context, windows.CERT_STORE_ADD_ALWAYS, nil)
		windows.CertFreeCertificateContext(context)
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := removeExactCertificate(store, certs[0]); err != nil {
		t.Fatal(err)
	}
	if err := removeExactCertificate(store, certs[0]); err != nil {
		t.Fatal("repeat removal:", err)
	}
	removed, err := findExactCertificate(store, certs[0])
	if err != nil || removed != nil {
		t.Fatal("target remains", err)
	}
	kept, err := findExactCertificate(store, certs[1])
	if err != nil || kept == nil {
		t.Fatal("other certificate removed", err)
	}
	windows.CertFreeCertificateContext(kept)
}

func TestCertificateStatusDoesNotCreateFiles(t *testing.T) {
	t.Chdir(t.TempDir())
	status, err := CaptureCertificateStatus()
	if err != nil || status != nil {
		t.Fatal(status, err)
	}
	entries, err := os.ReadDir(".")
	if err != nil || len(entries) != 0 {
		t.Fatal("status created files", err)
	}
}

func TestRemovalRejectsChangedReview(t *testing.T) {
	t.Chdir(t.TempDir())
	if _, err := PrepareCaptureCertificate(); err != nil {
		t.Fatal(err)
	}
	if err := RemoveReviewedCaptureCertificate("wrong thumbprint"); err == nil {
		t.Fatal("accepted wrong identity")
	}
}
