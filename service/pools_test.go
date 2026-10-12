package service

import (
	"ElmoBeacon/pb"
	"encoding/binary"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNormalizePoolTypes(t *testing.T) {
	got := NormalizePoolTypes([]int64{7, 6, 42, 0, -1, 6}, []int64{1, 42, 99})
	if !reflect.DeepEqual(got, []int64{1, 6, 7, 42, 99}) {
		t.Fatal(got)
	}
}

func gachaFixture(typ uint64, fields map[protowire.Number]string) []byte {
	row := protowire.AppendTag(nil, 1, protowire.VarintType)
	row = protowire.AppendVarint(row, 100)
	row = protowire.AppendTag(row, 4, protowire.VarintType)
	row = protowire.AppendVarint(row, typ)
	for field, text := range fields {
		row = protowire.AppendTag(row, field, protowire.BytesType)
		row = protowire.AppendString(row, text)
	}
	return protowire.AppendBytes(protowire.AppendTag(nil, 1, protowire.BytesType), row)
}

func TestGachaLayouts(t *testing.T) {
	for name, fields := range map[string]map[protowire.Number]string{
		"current": {21: "5:0.0189", 22: "5:1039:1015,4:1009", 27: "5:1039:1040", 35: "new artwork", 99: "future field"},
		"legacy":  {21: "5:1039:1015,4:1009", 29: "5:1039:1040"},
	} {
		t.Run(name, func(t *testing.T) {
			pools, ranks, err := parseGachaPools(gachaFixture(3, fields))
			if err != nil {
				t.Fatal(err)
			}
			if len(pools[100].UpItem) != 2 {
				t.Fatal(pools)
			}
			if _, ok := pools[100].UpItem[1015]; ok {
				t.Fatal("standard item counted as featured")
			}
			if ranks[1039] != 5 || ranks[1009] != 4 {
				t.Fatal(ranks)
			}
		})
	}
}

func TestSelectorDefaultIsNotHistoricalChoice(t *testing.T) {
	for _, typ := range []uint64{6, 7} {
		pools, _, err := parseGachaPools(gachaFixture(typ, map[protowire.Number]string{27: "5:1039", 29: "5:10393"}))
		if err != nil || len(pools[100].UpItem) != 0 {
			t.Fatalf("%v %v", pools, err)
		}
	}
}

func TestGachaUnknownTypeAndMalformedData(t *testing.T) {
	pools, _, err := parseGachaPools(gachaFixture(42, map[protowire.Number]string{99: "future"}))
	if err != nil || pools[100].PoolType != 42 {
		t.Fatalf("%v %v", pools, err)
	}
	for _, data := range [][]byte{{0xff}, {0x0a, 0x7f, 0x01}, {0x0a, 0x01, 0xff}} {
		if _, _, err := parseGachaPools(data); err == nil {
			t.Fatalf("accepted malformed data %x", data)
		}
	}
}

func TestRankItems(t *testing.T) {
	got := rankItems("5：1039：1040, 4: 1009,3:11022,5:0.5:bad:-1234")
	want := map[int64]int{1039: 5, 1040: 5, 1009: 4, 11022: 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatal(got)
	}
}

func TestPoolTypeServerDirectory(t *testing.T) {
	dir := t.TempDir()
	tableDir := filepath.Join(dir, "Table", "us")
	if err := os.MkdirAll(tableDir, 0755); err != nil {
		t.Fatal(err)
	}
	body, err := proto.Marshal(&pb.GachaTypeListData{Units: []*pb.GachaTypeList{{Id: 7}, {Id: 6}, {Id: 42}, {Id: 6}, {Id: 0}}})
	if err != nil {
		t.Fatal(err)
	}
	header := make([]byte, 4)
	binary.LittleEndian.PutUint32(header, 0)
	if err := os.WriteFile(filepath.Join(tableDir, "GachaTypeListData.bytes"), append(header, body...), 0600); err != nil {
		t.Fatal(err)
	}
	got, err := GetGachaPoolTypeList(dir, "us")
	if err != nil || !reflect.DeepEqual(got, []int64{6, 7, 42}) {
		t.Fatalf("%v %v", got, err)
	}
}

// Optional read-only check against a local client; no account or network required.
func TestInstalledGachaTables(t *testing.T) {
	dir := os.Getenv("ELMOBEACON_TEST_GAME_DATA")
	if dir == "" {
		t.Skip("set ELMOBEACON_TEST_GAME_DATA to check local client tables")
	}
	ids, err := GetGachaPoolTypeList(dir, "us")
	if err != nil || len(ids) == 0 {
		t.Fatalf("types: %v %v", ids, err)
	}
	pools, ranks, err := GetGachaPoolMap(dir, "us")
	if err != nil || len(pools) == 0 || len(ranks) == 0 {
		t.Fatalf("pools=%d ranks=%d err=%v", len(pools), len(ranks), err)
	}
	t.Logf("Parsed %d categories, %d banners, %d item ranks", len(ids), len(pools), len(ranks))
}

func TestLegacyWeaponLayout(t *testing.T) {
	pools, ranks, err := parseGachaPools(gachaFixture(4, map[protowire.Number]string{22: "5:10393:11016,4:11015", 28: "5:10393"}))
	if err != nil {
		t.Fatal(err)
	}
	if len(pools[100].UpItem) != 1 || ranks[11015] != 4 {
		t.Fatalf("%v %v", pools, ranks)
	}
	if _, ok := pools[100].UpItem[10393]; !ok {
		t.Fatal("missing featured weapon")
	}
}
