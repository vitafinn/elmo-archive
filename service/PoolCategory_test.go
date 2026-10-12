package service

import "testing"

func TestSeparateReunionFromCustom(t *testing.T) {
	pools := map[int64]GachaPoolInfo{999001: {PoolType: 6, NameID: 559138}, 999002: {PoolType: 7, NameID: 559139}}
	for _, tc := range []struct{ api, banner, display int64 }{
		{6, 252001, ReunionDolls}, {7, 253001, ReunionWeapons},
		{6, 1200011, 6}, {7, 1300011, 7}, {6, 999001, ReunionDolls},
		{7, 999002, ReunionWeapons}, {42, 252001, 42},
	} {
		got := DisplayPoolType(tc.api, tc.banner, pools)
		if got != tc.display || APIPoolType(got) != tc.api {
			t.Fatalf("%+v got %d", tc, got)
		}
	}
}
