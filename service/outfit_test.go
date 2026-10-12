package service

import (
	"google.golang.org/protobuf/encoding/protowire"
	"os"
	"testing"
)

func TestOutfitRewardsIncludeCurrencyAndLimitedItems(t *testing.T) {
	row := protowire.AppendVarint(protowire.AppendTag(nil, 1, protowire.VarintType), 98001)
	row = protowire.AppendVarint(protowire.AppendTag(row, 11, protowire.VarintType), 31083)
	row = protowire.AppendString(protowire.AppendTag(row, 7, protowire.BytesType), "5:31083,4:31086")
	row = protowire.AppendString(protowire.AppendTag(row, 9, protowire.BytesType), "4:110067,3:272:310054")
	data := protowire.AppendBytes(protowire.AppendTag(nil, 1, protowire.BytesType), row)
	pools, err := parseOutfitPools(data)
	if err != nil {
		t.Fatal(err)
	}
	p := pools[98001]
	if p.MainItem != 31083 || p.Rewards[31083] != "limited" || p.Rewards[31086] != "limited" || p.Rewards[110067] != "rare" || p.Rewards[272] != "normal" {
		t.Fatal(p)
	}
	if _, err := parseOutfitPools([]byte{10, 127}); err == nil {
		t.Fatal("accepted truncated data")
	}
}
func TestInstalledOutfitTables(t *testing.T) {
	dir := os.Getenv("ELMOBEACON_TEST_GAME_DATA")
	if dir == "" {
		t.Skip("no client path")
	}
	pools, err := GetOutfitPoolMap(dir, "us")
	if err != nil {
		t.Fatal(err)
	}
	if pools[98001].MainItem != 31083 || pools[98001].Rewards[272] != "normal" {
		t.Fatal(pools[98001])
	}
}
