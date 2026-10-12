package service

import "testing"

func TestCaptureOnlyDecryptsExplicitHosts(t *testing.T) {
	for _, host := range []string{"example.com:443", "gf2-gacha-record-us.sunborngame.com.attacker.test:443", "attacker@gf2-gacha-record-us.sunborngame.com", "gf2-gacha-record-us.sunborngame.com:80", "gf2-gacha-record-us.sunborngame.com:8443", "127.0.0.1:443", ""} {
		if shouldDecryptCaptureHost(host, "webshop") {
			t.Errorf("unexpected interception: %s", host)
		}
	}
	if !shouldDecryptCaptureHost("GF2-GACHA-RECORD-US.SUNBORNGAME.COM:443", "manual") {
		t.Fatal("record host rejected")
	}
	if shouldDecryptCaptureHost(captureWebshopHost, "manual") {
		t.Fatal("manual mode intercepted webshop")
	}
	if !shouldDecryptCaptureHost(captureWebshopHost+":443", "webshop") {
		t.Fatal("webshop host rejected")
	}
}
