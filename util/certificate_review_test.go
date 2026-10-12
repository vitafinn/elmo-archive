package util

import (
	"os"
	"strings"
	"testing"
)

func TestCertificateReviewKeepsTheSameLocalIdentity(t *testing.T) {
	t.Chdir(t.TempDir())
	first, err := PrepareCaptureCertificate()
	if err != nil {
		t.Fatal(err)
	}
	if first.Name != "Elmo Archive Local Capture CA" || first.Legacy || first.Trusted {
		t.Fatalf("unexpected new identity: %+v", first)
	}
	if len(strings.ReplaceAll(first.Thumbprint, " ", "")) != 40 {
		t.Fatal("expected Windows SHA-1 thumbprint")
	}
	before, err := os.ReadFile("ca.key")
	if err != nil {
		t.Fatal(err)
	}
	again, err := PrepareCaptureCertificate()
	if err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile("ca.key")
	if err != nil {
		t.Fatal(err)
	}
	if first.Thumbprint != again.Thumbprint || string(before) != string(after) {
		t.Fatal("review regenerated the certificate")
	}
	// A stale/mismatched review must fail without opening a Windows installer.
	if InstallReviewedCertificate("different thumbprint") == nil {
		t.Fatal("accepted an unreviewed certificate")
	}
	if certificatePromptActive.Load() {
		t.Fatal("opened native installer for mismatched certificate")
	}
}

func TestCertificateReviewDoesNotOverwriteBrokenPairs(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile("ca.crt", []byte("invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := PrepareCaptureCertificate(); err == nil {
		t.Fatal("accepted incomplete certificate pair")
	}
	data, _ := os.ReadFile("ca.crt")
	if string(data) != "invalid" {
		t.Fatal("overwrote existing certificate")
	}
}
