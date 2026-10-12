package handler

import (
	"ElmoBeacon/db"
	"ElmoBeacon/model"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

func TestPoolTypesRetainOnlySelectedUsersHistory(t *testing.T) {
	engine, err := xorm.NewEngine("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	engine.SetMaxOpenConns(1)
	engine.SetMapper(names.GonicMapper{})
	previous := db.Engine
	db.Engine = engine
	t.Cleanup(func() { db.Engine = previous; engine.Close() })
	if err := engine.Sync(new(model.User), new(model.Record)); err != nil {
		t.Fatal(err)
	}
	user := model.User{GameDataDir: t.TempDir(), GameServer: "us"}
	if _, err := engine.Insert(&user); err != nil {
		t.Fatal(err)
	}
	records := []model.Record{{Id: 1, UserId: user.Id, PoolType: 42}, {Id: 2, UserId: user.Id + 1, PoolType: 99}, {Id: 3, UserId: user.Id, PoolType: 6, PoolId: 252001}, {Id: 4, UserId: user.Id, PoolType: 7, PoolId: 253001}}
	if _, err := engine.Insert(&records); err != nil {
		t.Fatal(err)
	}
	got, err := NewApp().GetPoolTypes(user.Id)
	if err != nil || !reflect.DeepEqual(got, []int64{1, 3, 4, 5, 6, 7, 8, 9, 42, -6, -7}) {
		t.Fatalf("%v %v", got, err)
	}
	if _, err := NewApp().GetPoolTypes(user.Id + 100); err == nil {
		t.Fatal("expected missing-user error")
	}
}

func TestInstalledReunionHistory(t *testing.T) {
	path := os.Getenv("ELMOBEACON_TEST_DB")
	if path == "" {
		t.Skip("set ELMOBEACON_TEST_DB for a read-only history check")
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}
	engine, err := xorm.NewEngine("sqlite", filepath.ToSlash(absolute)+"?mode=ro")
	if err != nil {
		t.Fatal(err)
	}
	engine.SetMapper(names.GonicMapper{})
	previous := db.Engine
	db.Engine = engine
	t.Cleanup(func() { db.Engine = previous; engine.Close() })
	var record model.Record
	found, err := engine.Where("pool_id = ?", 252001).Get(&record)
	if err != nil || !found {
		t.Fatalf("reunion record not found: %v", err)
	}
	info, err := NewApp().GetPoolInfo(record.UserId, -6)
	if err != nil {
		t.Fatal(err)
	}
	if info.TotalCount != 80 || info.StoredCount+info.ExpiredPity+info.PendingPity != 20 || len(info.RecordList) != 1 || info.RecordList[0].Count != 60 || info.RecordList[0].IsMissing || info.RecordList[0].Outcome != "initial" || info.Rank5Avg != 60 || info.Rank5Spent != 80 || info.TargetAttempts != 0 {
		t.Fatalf("unexpected Reunion stats: %+v", info)
	}
	var oldRecords []model.Record
	if err := engine.Where("user_id = ? AND pool_type = 6 AND pool_id != 252001", record.UserId).Find(&oldRecords); err != nil {
		t.Fatal(err)
	}
	custom, err := NewApp().GetPoolInfo(record.UserId, 6)
	if err != nil {
		t.Fatal(err)
	}
	if custom.TotalCount != int64(len(oldRecords)) {
		t.Fatalf("custom %d, expected %d", custom.TotalCount, len(oldRecords))
	}
	t.Logf("Reunion: 80 total, Centaureissi at 60, 20 toward next five-star, no loss; Custom: %d separate pulls", custom.TotalCount)
}
