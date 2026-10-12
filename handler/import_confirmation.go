package handler

import (
	"ElmoBeacon/service"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strconv"
	"time"
)

type ImportPreview struct {
	Ticket         string
	UID            string
	Server         string
	IdentitySource string
}
type preparedImport struct {
	preview ImportPreview
	info    service.GameUserInfo
	expires time.Time
}

// Prepare reads connection information only. Secrets stay in the backend and
// no records or account entries are fetched/saved before confirmation.
func (a *App) PrepareImport(source string) (ImportPreview, error) {
	a.importMu.Lock()
	defer a.importMu.Unlock()
	if a.importBusy {
		return ImportPreview{}, fmt.Errorf("import already running")
	}
	a.pendingImport = nil
	var info *service.GameUserInfo
	var err error
	switch source {
	case "os":
		info, err = service.GetUserInfoOS()
	default:
		return ImportPreview{}, fmt.Errorf("invalid import source")
	}
	if err != nil {
		return ImportPreview{}, fmt.Errorf("could not read the connected account; reconnect its history and try again")
	}
	if info.Uid == 0 {
		return ImportPreview{}, fmt.Errorf("account identifier is missing; reconnect")
	}
	if info.GameServer == service.GameServerCN {
		return ImportPreview{}, fmt.Errorf("the China client is not supported by Elmo Archive")
	}
	if info.IdentitySource == "webshop" && string(info.GameServer) != "us" {
		return ImportPreview{}, fmt.Errorf("Darkwinter webshop cannot identify a Haoplay account; reconnect using its game UID")
	}
	bytes := make([]byte, 24)
	if _, err = rand.Read(bytes); err != nil {
		return ImportPreview{}, err
	}
	method := info.IdentitySource
	switch method {
	case "manual", "webshop", "hash", "gameLog":
	default:
		method = "saved"
	}
	preview := ImportPreview{Ticket: hex.EncodeToString(bytes), UID: strconv.FormatUint(info.Uid, 10), Server: string(info.GameServer), IdentitySource: method}
	a.pendingImport = &preparedImport{preview: preview, info: *info, expires: time.Now().Add(5 * time.Minute)}
	return preview, nil
}

func (a *App) takePreparedImport(ticket, uid string) (*service.GameUserInfo, error) {
	a.importMu.Lock()
	defer a.importMu.Unlock()
	pending := a.pendingImport
	if a.importBusy || pending == nil || pending.preview.Ticket != ticket || pending.preview.UID != uid || time.Now().After(pending.expires) {
		return nil, fmt.Errorf("account confirmation expired or changed; review the account again")
	}
	a.pendingImport = nil
	a.importBusy = true
	info := pending.info
	return &info, nil
}

func (a *App) ConfirmImport(ticket, uid string) (*service.SyncResult, error) {
	info, err := a.takePreparedImport(ticket, uid)
	if err != nil {
		return nil, err
	}
	defer func() { a.importMu.Lock(); a.importBusy = false; a.importMu.Unlock() }()
	return service.SyncRecords(info, func(event service.ImportProgress) {
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "importProgress", event)
		}
	})
}
