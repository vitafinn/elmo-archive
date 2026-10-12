package service

import "testing"

func TestChinaSyncIsUnsupported(t *testing.T) {
	called := false
	_, err := SyncRecords(&GameUserInfo{GameServer: GameServerCN}, func(ImportProgress) { called = true })
	if err == nil || called {
		t.Fatal("China sync must stop before progress, database access or network requests")
	}
}
