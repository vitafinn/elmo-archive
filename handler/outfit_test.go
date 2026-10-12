package handler

import (
	"ElmoBeacon/db"
	"ElmoBeacon/model"
	"ElmoBeacon/service"
	"os"
	"path/filepath"
	"testing"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

func TestOutfitSeasonIsolationAndUnknownRewards(t *testing.T) {
	records := []model.Record{{Id: 1, PoolType: 9, PoolId: 1, ItemId: 272}, {Id: 2, PoolType: 9, PoolId: 2, ItemId: 999999}, {Id: 3, PoolType: 9, PoolId: 1, ItemId: 31083}, {Id: 4, PoolType: 10, PoolId: 1, ItemId: 31083}}
	pools := map[int64]service.GachaPoolInfo{1: {PoolType: 9, StartTime: 1, EndTime: 100}, 2: {PoolType: 9, StartTime: 1, EndTime: 200}, 3: {PoolType: 9, StartTime: 1, EndTime: 200}}
	outfits := map[int64]service.OutfitPoolInfo{1: {MainItem: 31083, Rewards: map[int64]string{272: "normal", 31083: "limited"}}}
	info := buildOutfitInfo(records, 9, pools, outfits, nil, nil, 150)
	if info.TotalCount != 3 || len(info.Seasons) != 3 {
		t.Fatal(info)
	}
	for _, s := range info.Seasons {
		if s.ID == 1 && (s.TotalCount != 2 || s.MainOutfitAt != 2 || s.NormalCount != 1 || s.LimitedCount != 1 || s.Status != "expired" || s.Records[0].AccessNumber != 2) {
			t.Fatal(s)
		}
		if s.ID == 2 && (s.UnknownCount != 1 || s.Status != "active" || s.Records[0].Name != "Item #999999") {
			t.Fatal(s)
		}
	}
	unknown := buildOutfitInfo(records[:1], 9, nil, nil, nil, nil, 150)
	if unknown.TotalCount != 1 || unknown.Seasons[0].Status != "unknown" || unknown.Seasons[0].Records[0].Category != "unknown" {
		t.Fatal(unknown)
	}
}
func TestInstalledOutfitHistory(t *testing.T) {
	path := os.Getenv("ELMOBEACON_TEST_DB")
	if path == "" {
		t.Skip("no database path")
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
	old := db.Engine
	db.Engine = engine
	t.Cleanup(func() { db.Engine = old; engine.Close() })
	var first model.Record
	found, err := engine.Where("pool_type = 9").Get(&first)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Skip("no outfit records yet")
	}
	var records []model.Record
	if err := engine.Where("pool_type = 9 AND user_id = ?", first.UserId).Find(&records); err != nil {
		t.Fatal(err)
	}
	info, err := NewApp().GetOutfitInfo(first.UserId, 9)
	if err != nil {
		t.Fatal(err)
	}
	var total int64
	for _, s := range info.Seasons {
		total += s.TotalCount
		if int64(len(s.Records)) != s.TotalCount || s.LimitedCount+s.RareCount+s.NormalCount+s.UnknownCount != s.TotalCount {
			t.Fatal(s)
		}
	}
	if info.TotalCount != int64(len(records)) || total != info.TotalCount {
		t.Fatal(info)
	}
	t.Logf("Verified %d imported accesses across %d displayed seasons", total, len(info.Seasons))
}
