package handler

import "testing"

func TestExpiryPreservesSpendingAndCompletedAverage(t *testing.T) {
	for _, now := range []int64{150, 201} {
		info := PoolInfo{TotalCount: 80, Rank5Count: 1}
		periods := map[int64]*BannerPeriod{1: {ID: 1, StartTime: 100, EndTime: 200, TotalCount: 80, Rank5Count: 1, CompletedPulls: 60, TailCount: 20}}
		finalizeSelectorPeriods(&info, periods, now)
		p := info.BannerPeriods[0]
		if p.Rank5Avg != 60 || p.Rank5Spent != 80 || p.TailCount != 20 {
			t.Fatal(p)
		}
		if now == 150 && (info.StoredCount != 20 || info.ExpiredPity != 0) {
			t.Fatal(info)
		}
		if now == 201 && (info.StoredCount != 0 || info.ExpiredPity != 20) {
			t.Fatal(info)
		}
	}
}
func TestUnknownDatesAndZeroDrops(t *testing.T) {
	info := PoolInfo{}
	periods := map[int64]*BannerPeriod{1: {ID: 1, TotalCount: 20, TailCount: 20}, 2: {ID: 2, StartTime: 1, EndTime: 2, TotalCount: 7, TailCount: 7}}
	finalizeSelectorPeriods(&info, periods, 10)
	if info.PendingPity != 20 || info.ExpiredPity != 7 || info.StoredCount != 0 {
		t.Fatal(info)
	}
	for _, p := range info.BannerPeriods {
		if p.Rank5Avg != 0 || p.Rank5Spent != 0 {
			t.Fatal(p)
		}
	}
}

func TestAggregateWeightsDropsAndIncludesEmptyBanners(t *testing.T) {
	info := PoolInfo{}
	periods := map[int64]*BannerPeriod{
		1: {ID: 1, TotalCount: 80, Rank5Count: 1, CompletedPulls: 60, TailCount: 20},
		2: {ID: 2, TotalCount: 100, Rank5Count: 2, CompletedPulls: 90, TailCount: 10},
		3: {ID: 3, TotalCount: 20, TailCount: 20},
	}
	finalizeSelectorPeriods(&info, periods, 10)
	if info.Rank5Avg != 50 || info.Rank5Spent != 200.0/3 {
		t.Fatal(info)
	}
}
