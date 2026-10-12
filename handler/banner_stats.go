package handler

import (
	"ElmoBeacon/service"
	"slices"
)

type BannerPeriod struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	StartTime        int64   `json:"startTime"`
	EndTime          int64   `json:"endTime"`
	LastTimestamp    int64   `json:"lastTimestamp"`
	Status           string  `json:"status"`
	TotalCount       int64   `json:"totalCount"`
	Rank5Count       int64   `json:"rank5Count"`
	CompletedPulls   int64   `json:"completedPulls"`
	Rank5Avg         float64 `json:"rank5Avg"`
	Rank5Spent       float64 `json:"rank5Spent"`
	TailCount        int64   `json:"tailCount"`
	TargetAttempts   int64   `json:"targetAttempts"`
	TargetWins       int64   `json:"targetWins"`
	UnknownOutcomes  int64   `json:"unknownOutcomes"`
	UnknownRankCount int64   `json:"unknownRankCount"`
}

func newBannerPeriod(id int64, info service.GachaPoolInfo) *BannerPeriod {
	return &BannerPeriod{ID: id, StartTime: info.StartTime, EndTime: info.EndTime}
}
func finalizeSelectorPeriods(info *PoolInfo, periods map[int64]*BannerPeriod, now int64) {
	info.StoredCount = 0
	var total, drops, completed int64
	for _, period := range periods {
		total += period.TotalCount
		drops += period.Rank5Count
		completed += period.CompletedPulls
		period.Status = "unknown"
		switch {
		case period.EndTime > 0 && now > period.EndTime:
			period.Status = "expired"
			info.ExpiredPity += period.TailCount
		case period.StartTime > 0 && period.StartTime <= now && now <= period.EndTime:
			period.Status = "active"
			info.StoredCount += period.TailCount
		default:
			info.PendingPity += period.TailCount
		}
		if period.Rank5Count > 0 {
			period.Rank5Avg = float64(period.CompletedPulls) / float64(period.Rank5Count)
			period.Rank5Spent = float64(period.TotalCount) / float64(period.Rank5Count)
		}
		info.BannerPeriods = append(info.BannerPeriods, *period)
	}
	if drops > 0 {
		info.Rank5Avg = float64(completed) / float64(drops)
		info.Rank5Spent = float64(total) / float64(drops)
	}
	slices.SortFunc(info.BannerPeriods, func(a, b BannerPeriod) int {
		if a.StartTime > b.StartTime {
			return -1
		}
		if a.StartTime < b.StartTime {
			return 1
		}
		if a.ID > b.ID {
			return -1
		}
		if a.ID < b.ID {
			return 1
		}
		return 0
	})
}
