package handler

import (
	"ElmoBeacon/db"
	"ElmoBeacon/model"
	"os"
	"path/filepath"
	"testing"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

func TestInstalledCharacterHistory(t *testing.T) {
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
	found, err := engine.Where("item_id = ?", 1028).Get(&first)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Skip("no Centaureissi records")
	}
	expected, err := engine.Where("user_id = ? AND item_id = ?", first.UserId, first.ItemId).Count(new(model.Record))
	if err != nil {
		t.Fatal(err)
	}
	records, err := NewApp().GetItemHistory(first.UserId, first.ItemId)
	if err != nil {
		t.Fatal(err)
	}
	if int64(len(records)) != expected {
		t.Fatalf("got %d copies; raw history has %d", len(records), expected)
	}
	seen := map[string]bool{}
	for i, record := range records {
		if record.ItemID != first.ItemId || record.ItemType != 10 || seen[record.ID] {
			t.Fatalf("incorrect identity or duplicated record: %+v", record)
		}
		seen[record.ID] = true
		if i > 0 && records[i-1].Timestamp > record.Timestamp {
			t.Fatal("history is not chronological")
		}
		if record.Count <= 0 {
			t.Fatalf("missing acquisition sequence: %d", record.Count)
		}
	}
	t.Logf("Verified %d copies across imported categories without duplicating Custom/Reunion", len(records))
	weapon, err := NewApp().GetDollWeaponHistory(first.UserId, first.ItemId)
	if err != nil {
		t.Fatal(err)
	}
	if weapon.ItemID != 11056 || weapon.Name == "" {
		t.Fatalf("incorrect Centaureissi weapon: %d %q", weapon.ItemID, weapon.Name)
	}
	weaponCount, err := engine.Where("user_id = ? AND item_id = ?", first.UserId, weapon.ItemID).Count(new(model.Record))
	if err != nil {
		t.Fatal(err)
	}
	if int64(len(weapon.Records)) != weaponCount {
		t.Fatal("weapon history count differs from raw records")
	}
	for i, record := range weapon.Records {
		if record.ItemID != weapon.ItemID {
			t.Fatal("wrong weapon record")
		}
		if i > 0 && compareHistoryRecords(weapon.Records[i-1], record) > 0 {
			t.Fatal("weapon history not chronological")
		}
	}
}
