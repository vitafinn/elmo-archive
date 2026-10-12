package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"hash/fnv"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"ElmoBeacon/logic"
	"ElmoBeacon/util"

	"github.com/elazarl/goproxy"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SessionConfig struct {
	IdentitySource     string
	FullGachaRecordUrl string
	AccessToken        string
	UidString          string
}

var (
	ActiveSession     SessionConfig
	sessionMutex      sync.Mutex
	captureMutex      sync.Mutex
	tcpListener       *net.Listener
	mitmServer        *http.Server
	captureEmitted    bool
	systemProxyActive bool
)

func setGoproxyCA() error {
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		return err
	}
	caKey, err := os.ReadFile("ca.key")
	if err != nil {
		return err
	}

	tlsCert, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	tlsCert.Leaf, err = x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return err
	}

	goproxy.GoproxyCa = tlsCert
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&tlsCert)}

	return nil
}

func generateHash(emailStr, accessToken string) string {
	identifier := emailStr
	authParts := strings.Split(accessToken, ".")
	if len(authParts) > 0 {
		padded := authParts[0]
		if m := len(padded) % 4; m != 0 {
			padded += strings.Repeat("=", 4-m)
		}
		if decoded, err := base64.StdEncoding.DecodeString(padded); err == nil {
			if m := regexp.MustCompile(`"openid":"([^"]+)"`).FindStringSubmatch(string(decoded)); len(m) > 1 {
				identifier = m[1]
			}
		}
	}
	if identifier != "" {
		h := fnv.New64a()
		h.Write([]byte(identifier))
		return fmt.Sprintf("%d", h.Sum64()&0x7FFFFFFFFFFFFFFF)
	}
	return "0"
}

func checkAndEmit(ctx context.Context, uidMode string) {
	if ActiveSession.AccessToken != "" && ActiveSession.FullGachaRecordUrl != "" {
		if uidMode == "webshop" && ActiveSession.UidString == "" {
			return
		}
		if !captureEmitted {
			captureEmitted = true
			log.Info().Msg("All session parameters locked in RAM. Triggering frontend UI.")
			runtime.EventsEmit(ctx, "captureSuccess")
		}
	}
}

func StartCapture(ctx context.Context, uidMode string) error {
	captureMutex.Lock()
	defer captureMutex.Unlock()
	if mitmServer != nil || systemProxyActive {
		return errors.New("capture engine is already running")
	}

	sessionMutex.Lock()
	method := uidMode
	if method != "webshop" && method != "hash" {
		method = "manual"
	}
	ActiveSession = SessionConfig{IdentitySource: method}
	captureEmitted = false
	sessionMutex.Unlock()

	if !util.IsTrustedCA() {
		return errors.New("review and install the local capture certificate before connecting")
	}

	if err := setGoproxyCA(); err != nil {
		return errors.Wrap(err, "failed to load CA for proxy")
	}

	port := 8899
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return errors.Wrap(err, "failed to start tcp listener")
	}

	proxy := goproxy.NewProxyHttpServer()
	// goproxy's default transport skips verification. Verify real server TLS.
	proxy.Tr = &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}}
	proxy.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		if shouldDecryptCaptureHost(host, uidMode) {
			return goproxy.MitmConnect, host
		}
		return goproxy.OkConnect, host // Unrelated HTTPS remains an encrypted tunnel.
	})

	proxy.OnRequest().DoFunc(
		func(req *http.Request, proxyCtx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if !shouldDecryptCaptureHost(req.Host, uidMode) || (req.URL.Scheme != "https" && req.TLS == nil) {
				return req, nil
			}
			scheme := "https"
			fullUrl := scheme + "://" + req.Host + req.URL.Path
			if req.URL.RawQuery != "" {
				fullUrl += "?" + req.URL.RawQuery
			}

			sessionMutex.Lock()
			currentUid := ActiveSession.UidString
			currentAccess := ActiveSession.AccessToken
			sessionMutex.Unlock()

			// Route 1: Webshop Node
			if uidMode == "webshop" && captureHostname(req.Host) == captureWebshopHost && currentUid == "" {
				webshopToken := req.Header.Get("Authorization")
				if webshopToken != "" && strings.HasPrefix(webshopToken, "eyJ") {
					uid, err := logic.FetchTrueUidFromWebshop(webshopToken)
					if err == nil && uid != "" {
						sessionMutex.Lock()
						ActiveSession.UidString = uid
						runtime.EventsEmit(ctx, "importProgress", ImportProgress{Stage: "webshopCaptured"})
						log.Info().Msg("Webshop account identifier captured")
						checkAndEmit(ctx, uidMode)
						sessionMutex.Unlock()
					}
				}
			}

			// Route 2 & 3: Game Node
			if isRecordCaptureHost(req.Host) && strings.HasSuffix(req.URL.Path, "/list") && req.Method == "POST" && currentAccess == "" {
				gameToken := req.Header.Get("Authorization")
				if gameToken != "" {
					sessionMutex.Lock()
					ActiveSession.FullGachaRecordUrl = fullUrl
					ActiveSession.AccessToken = gameToken
					runtime.EventsEmit(ctx, "importProgress", ImportProgress{Stage: "gameCaptured"})

					if uidMode != "webshop" {
						if uidMode == "hash" {
							emailStr := req.URL.Query().Get("u")
							ActiveSession.UidString = generateHash(emailStr, gameToken)
							log.Info().Msg("Node 2 Complete: Hash generated for Privacy Route")
						} else {
							ActiveSession.UidString = uidMode
							log.Info().Msg("Node 2 Complete: Manual UID injected")
						}
					} else {
						log.Info().Msg("Node 2 Complete: Game Token and URL locked. Waiting for Webshop UID...")
					}

					checkAndEmit(ctx, uidMode)
					sessionMutex.Unlock()
				}
			}
			return req, nil
		})

	tcpListener = &listener
	server := &http.Server{Handler: proxy}
	mitmServer = server

	go func() {
		server.Serve(listener)
		captureMutex.Lock()
		if mitmServer == server {
			if systemProxyActive {
				if err := util.DisableSysProxy(); err == nil {
					systemProxyActive = false
				}
			}
			tcpListener = nil
			mitmServer = nil
		}
		captureMutex.Unlock()
	}()

	if err = util.EnableSysProxy(port); err != nil {
		server.Close()
		mitmServer = nil
		tcpListener = nil
		return errors.Wrap(err, "failed to inject system proxy parameters")
	}
	systemProxyActive = true

	if uidMode == "webshop" {
		log.Info().Msg("Proxy running. Launching browser for Webshop authentication...")
		exec.Command("cmd", "/c", "start", "https://gf2-webshop-us.sunborngame.com/").Start()
	}

	return nil
}

func StopCapture() error {
	captureMutex.Lock()
	defer captureMutex.Unlock()
	return stopCaptureLocked()
}

func stopCaptureLocked() error {
	if systemProxyActive {
		if err := util.DisableSysProxy(); err != nil {
			return err
		}
		systemProxyActive = false
	}
	if mitmServer != nil {
		mitmServer.Close()
		mitmServer = nil
		tcpListener = nil
	}
	return nil
}
