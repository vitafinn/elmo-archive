package handler

// AcquisitionLoss exposes the preceding loss instead of hiding it in a total.
type AcquisitionLoss struct {
	ID         string
	Name       string
	Count      int64
	BannerID   int64
	BannerName string
}

// records must be chronological and belong to one display category. Selector
// guarantees are isolated by banner; ordinary limited guarantees can carry over.
func withAcquisitionCosts(records []DisplayRecord) []DisplayRecord {
	result := make([]DisplayRecord, len(records))
	previous := make(map[int64]DisplayRecord)
	for i, record := range records {
		record.AcquisitionCount = record.Count
		record.PreviousLoss = nil
		record.CostBasis = "direct"
		selector := record.PoolType == 6 || record.PoolType == 7 || record.PoolType == -6 || record.PoolType == -7
		limited := record.PoolType == 3 || record.PoolType == 4
		key := int64(0)
		if selector {
			key = record.BannerID
		}
		prev, hasPrevious := previous[key]
		if selector && (record.Outcome == "unknown" || record.Outcome == "") {
			record.CostBasis = "unknown"
		}
		if selector && record.Outcome == "guaranteed" && (!hasPrevious || prev.Outcome != "loss") {
			record.CostBasis = "incomplete"
		}
		link := false
		if hasPrevious && prev.ItemID != record.ItemID {
			if selector && prev.Outcome == "loss" && record.Outcome == "guaranteed" {
				link = true
				record.CostBasis = "recordedLoss"
			} else if limited && prev.IsMissing && !record.IsMissing {
				link = true
				switch {
				case !prev.TargetKnown || !record.TargetKnown:
					record.CostBasis = "estimatedOutcome"
				case prev.BannerID != record.BannerID:
					record.CostBasis = "carried"
				default:
					record.CostBasis = "linked"
				}
			}
		}
		if link {
			record.AcquisitionCount += prev.Count
			record.PreviousLoss = &AcquisitionLoss{ID: prev.ID, Name: prev.Name, Count: prev.Count, BannerID: prev.BannerID, BannerName: prev.BannerName}
		}
		result[i] = record
		previous[key] = record
	}
	return result
}
