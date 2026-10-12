package service

import (
	"os"
	"testing"
)

func TestInspectBannerEvidence(t *testing.T) {
	dir := os.Getenv("ELMOBEACON_TEST_GAME_DATA")
	if dir == "" {
		t.Skip()
	}
	pools, _, err := GetGachaPoolMap(dir, "us")
	if err != nil {
		t.Fatal(err)
	}
	lang, err := GetLangDataMap(dir, "us", "en")
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []int64{2000011, 252001, 253001} {
		p, ok := pools[id]
		t.Logf("banner=%d inCurrentTable=%v apiType=%d nameID=%d name=%q displayType=%d", id, ok, p.PoolType, p.NameID, lang[p.NameID], DisplayPoolType(p.PoolType, id, pools))
	}
}
