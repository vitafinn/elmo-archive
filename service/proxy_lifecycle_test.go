package service

import "testing"

// With no capture active, stopping must not alter Windows proxy settings.
func TestStopCaptureWithoutSession(t *testing.T) {
	captureMutex.Lock()
	if mitmServer != nil || systemProxyActive {
		captureMutex.Unlock()
		t.Fatal("unexpected active capture")
	}
	captureMutex.Unlock()
	if err := StopCapture(); err != nil {
		t.Fatal(err)
	}
	if err := StopCapture(); err != nil {
		t.Fatal(err)
	}
}
