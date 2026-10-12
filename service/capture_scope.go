package service

import (
	"net"
	"strings"
)

const captureWebshopHost = "gf2-us-webshop-api.sunborngame.com"

// Exact hostnames only: no wildcard suffix or substring matching. Haoplay paths
// remain preview compatibility entries pending testing on those clients.
var captureRecordHosts = []string{
	"gf2-gacha-record-us.sunborngame.com",
	"gf2-gacha-record-intl.sunborngame.com", "gf2-gacha-record-jp.sunborngame.com",
	"gf2-gacha-record-kr.sunborngame.com", "gf2-gacha-record-asia.sunborngame.com",
	"gf2-gacha-record-intl.haoplay.com", "gf2-gacha-record-jp.haoplay.com",
	"gf2-gacha-record-kr.haoplay.com", "gf2-gacha-record-asia.haoplay.com",
}

func captureHostname(authority string) string {
	host := authority
	if strings.Contains(authority, ":") {
		var port string
		var err error
		host, port, err = net.SplitHostPort(authority)
		if err != nil || port != "443" {
			return ""
		}
	}
	return strings.TrimSuffix(strings.ToLower(host), ".")
}

func isRecordCaptureHost(authority string) bool {
	host := captureHostname(authority)
	for _, allowed := range captureRecordHosts {
		if host == allowed {
			return true
		}
	}
	return false
}

func shouldDecryptCaptureHost(authority, method string) bool {
	return isRecordCaptureHost(authority) || (method == "webshop" && captureHostname(authority) == captureWebshopHost)
}
