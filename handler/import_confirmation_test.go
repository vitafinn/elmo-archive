package handler

import (
	"ElmoBeacon/service"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func pendingFixture() *App {
	return &App{pendingImport: &preparedImport{preview: ImportPreview{Ticket: "ticket", UID: "123", Server: "us", IdentitySource: "webshop"}, info: service.GameUserInfo{Uid: 123, GameServer: service.GameServerUS, GameAccessToken: "secret", GachaRecordUrl: "private-url"}, expires: time.Now().Add(time.Minute)}}
}
func TestImportConfirmation(t *testing.T) {
	app := pendingFixture()
	if _, err := app.takePreparedImport("wrong", "123"); err == nil {
		t.Fatal("accepted stale ticket")
	}
	if _, err := app.takePreparedImport("ticket", "456"); err == nil {
		t.Fatal("accepted different UID")
	}
	info, err := app.takePreparedImport("ticket", "123")
	if err != nil || info.Uid != 123 || info.GameAccessToken != "secret" {
		t.Fatal("lost confirmed snapshot")
	}
	if _, err = app.takePreparedImport("ticket", "123"); err == nil {
		t.Fatal("accepted replay")
	}
	app = pendingFixture()
	app.pendingImport.expires = time.Now().Add(-time.Second)
	if _, err = app.takePreparedImport("ticket", "123"); err == nil {
		t.Fatal("accepted expired confirmation")
	}
	app = pendingFixture()
	app.importBusy = true
	if _, err = app.takePreparedImport("ticket", "123"); err == nil {
		t.Fatal("accepted concurrent import")
	}
}
func TestImportPreviewContainsNoCredentials(t *testing.T) {
	bytes, err := json.Marshal(pendingFixture().pendingImport.preview)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(bytes), "secret") || strings.Contains(string(bytes), "private-url") {
		t.Fatal("preview exposed credentials")
	}
}
func TestChinaImportIsNotOffered(t *testing.T) {
	app := pendingFixture()
	if _, err := app.PrepareImport("cn"); err == nil {
		t.Fatal("China imports must be rejected before reading client data")
	}
	if app.pendingImport != nil {
		t.Fatal("old confirmation remained active")
	}
}
