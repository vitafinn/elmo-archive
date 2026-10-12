package logic

import "testing"

func TestWebshopUIDPrecision(t *testing.T) {
	for _, body := range []string{`{"data":{"uid":9007199254740993}}`, `{"uid":"9007199254740993"}`} {
		uid, err := parseWebshopUID([]byte(body))
		if err != nil || uid != "9007199254740993" {
			t.Fatalf("uid=%s err=%v", uid, err)
		}
	}
	for _, body := range []string{`{"uid":0}`, `{"uid":"someone@example.com"}`, `{"data":{}}`, `{"uid":1.5}`} {
		if _, err := parseWebshopUID([]byte(body)); err == nil {
			t.Fatal("accepted invalid UID")
		}
	}
}
