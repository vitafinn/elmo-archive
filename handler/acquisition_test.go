package handler

import "testing"

func TestAcquisitionCosts(t *testing.T) {
	makeRecord := func(id string, item, pool, banner, count int64, outcome string, missing, known bool) DisplayRecord {
		return DisplayRecord{ID: id, ItemID: item, PoolType: pool, BannerID: banner, Count: count, Outcome: outcome, IsMissing: missing, TargetKnown: known}
	}
	tests := []struct {
		name   string
		rows   []DisplayRecord
		index  int
		cost   int64
		basis  string
		linked bool
	}{
		{"limited loss plus target", []DisplayRecord{makeRecord("1", 1, 3, 10, 65, "", true, true), makeRecord("2", 2, 3, 10, 60, "", false, true)}, 1, 125, "linked", true},
		{"selector recorded loss sequence", []DisplayRecord{makeRecord("1", 1, -6, 10, 65, "loss", true, false), makeRecord("2", 2, -6, 10, 60, "guaranteed", false, false)}, 1, 125, "recordedLoss", true},
		{"selector expiry isolates", []DisplayRecord{makeRecord("1", 1, -6, 10, 65, "loss", true, false), makeRecord("2", 2, -6, 11, 60, "initial", false, false)}, 1, 60, "direct", false},
		{"unknown is not a loss", []DisplayRecord{makeRecord("1", 1, -6, 10, 65, "unknown", false, false), makeRecord("2", 2, -6, 10, 60, "unknown", false, false)}, 1, 60, "unknown", false},
		{"limited carryover", []DisplayRecord{makeRecord("1", 1, 3, 10, 65, "", true, true), makeRecord("2", 2, 3, 11, 60, "", false, true)}, 1, 125, "carried", true},
		{"missing metadata", []DisplayRecord{makeRecord("1", 1, 3, 10, 65, "", true, false), makeRecord("2", 2, 3, 10, 60, "", false, true)}, 1, 125, "estimatedOutcome", true},
		{"same character not double counted", []DisplayRecord{makeRecord("1", 2, 3, 10, 65, "", true, true), makeRecord("2", 2, 3, 11, 60, "", false, true)}, 1, 60, "direct", false},
		{"no repeated loss", []DisplayRecord{makeRecord("1", 1, 3, 10, 65, "", true, true), makeRecord("2", 2, 3, 10, 60, "", false, true), makeRecord("3", 2, 3, 10, 10, "", false, true)}, 2, 10, "direct", false},
		{"incomplete guarantee", []DisplayRecord{makeRecord("1", 2, 6, 10, 60, "guaranteed", false, false)}, 0, 60, "incomplete", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := withAcquisitionCosts(tt.rows)[tt.index]
			if got.AcquisitionCount != tt.cost || got.CostBasis != tt.basis || (got.PreviousLoss != nil) != tt.linked {
				t.Fatalf("%+v", got)
			}
		})
	}
}
