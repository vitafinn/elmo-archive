package handler

import (
	"ElmoBeacon/pb"
	"reflect"
	"testing"
)

func TestSearchNamesUseSharedItemKey(t *testing.T) {
	items := map[int64]*pb.Item{1: {Name: &pb.LanguageStringData{Id: 42}}, 2: {}}
	languages := []map[int64]string{{42: "瓊玖"}, {42: "Qiongjiu", 1: "Wrong item"}, {42: "瓊玖"}, nil}
	if got := itemSearchNames(1, items, languages); !reflect.DeepEqual(got, []string{"瓊玖", "Qiongjiu"}) {
		t.Fatal(got)
	}
	for _, id := range []int64{2, 999} {
		if got := itemSearchNames(id, items, languages); len(got) != 0 {
			t.Fatal(got)
		}
	}
}

func TestSearchMissingLanguageFilesKeepCurrentNames(t *testing.T) {
	current := map[int64]string{42: "Qiongjiu"}
	got := searchLanguages(t.TempDir(), "us", current)
	if len(got) != 1 || got[0][42] != "Qiongjiu" {
		t.Fatal(got)
	}
}
