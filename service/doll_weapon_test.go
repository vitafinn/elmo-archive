package service

import (
	"google.golang.org/protobuf/encoding/protowire"
	"os"
	"testing"
)

func weaponRow(doll, weapon uint64) []byte {
	row := protowire.AppendTag(nil, 1, 0)
	row = protowire.AppendVarint(row, doll)
	row = protowire.AppendTag(row, 24, 0)
	row = protowire.AppendVarint(row, weapon)
	return protowire.AppendBytes(protowire.AppendTag(nil, 1, 2), row)
}
func TestDollWeaponLinks(t *testing.T) {
	b := append(weaponRow(1028, 11056), weaponRow(1001, 0)...)
	links, err := parseDollWeapons(b)
	if err != nil || len(links) != 1 || links[1028] != 11056 {
		t.Fatalf("%v %v", links, err)
	}
	if _, err = parseDollWeapons(append(b, weaponRow(1028, 10393)...)); err == nil {
		t.Fatal("accepted conflicting links")
	}
	if _, err = parseDollWeapons([]byte{10, 255}); err == nil {
		t.Fatal("accepted truncated data")
	}
}
func TestInstalledDollWeaponLinks(t *testing.T) {
	dir := os.Getenv("ELMOBEACON_TEST_GAME_DATA")
	if dir == "" {
		t.Skip("no game tables")
	}
	links, err := GetDollWeaponMap(dir, "us")
	if err != nil {
		t.Fatal(err)
	}
	for doll, weapon := range map[int64]int64{1028: 11056, 1039: 10393, 1025: 11038} {
		if links[doll] != weapon {
			t.Fatalf("doll %d: got %d want %d", doll, links[doll], weapon)
		}
	}
	if links[1001] != 0 {
		t.Fatal("general recommendation mistaken for signature")
	}
	t.Logf("Validated %d elite weapon associations", len(links))
}
