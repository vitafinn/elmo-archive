package handler

import (
	"ElmoBeacon/db"
	"ElmoBeacon/model"
	"testing"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

func TestSelectorAnnotationIsolationAndReset(t *testing.T) {
	engine, err := xorm.NewEngine("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	engine.SetMaxOpenConns(1)
	engine.SetMapper(names.GonicMapper{})
	old := db.Engine
	db.Engine = engine
	t.Cleanup(func() { db.Engine = old; engine.Close() })
	if err := engine.Sync(new(model.Record), new(model.Setting)); err != nil {
		t.Fatal(err)
	}
	records := []model.Record{{Id: 10, UserId: 1, PoolType: 6, PoolId: 252001, ItemId: 1028}, {Id: 11, UserId: 1, PoolType: 3}}
	if _, err := engine.Insert(&records); err != nil {
		t.Fatal(err)
	}
	app := NewApp()
	if err := app.SetSelectorOutcome(1, "10", "loss"); err != nil {
		t.Fatal(err)
	}
	saved, err := app.selectorAnnotations(1)
	if err != nil || saved["10"] != "loss" {
		t.Fatalf("%v %v", saved, err)
	}
	if err := app.SetSelectorOutcome(2, "10", "win"); err == nil {
		t.Fatal("allowed another user")
	}
	if err := app.SetSelectorOutcome(1, "11", "win"); err == nil {
		t.Fatal("allowed non-selector")
	}
	if err := app.SetSelectorOutcome(1, "10", "invalid"); err == nil {
		t.Fatal("allowed invalid outcome")
	}
	if err := app.SetSelectorOutcome(1, "10", "auto"); err != nil {
		t.Fatal(err)
	}
	saved, err = app.selectorAnnotations(1)
	if err != nil || saved["10"] != "auto" {
		t.Fatalf("%v %v", saved, err)
	}
	var record model.Record
	if _, err := engine.Where("id = ?", 10).Get(&record); err != nil {
		t.Fatal(err)
	}
	if record.ItemId != 1028 || record.PoolId != 252001 {
		t.Fatal("annotation changed original history")
	}
}
