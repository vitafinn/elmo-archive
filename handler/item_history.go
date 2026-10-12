package handler

import (
	"ElmoBeacon/db"
	"ElmoBeacon/model"
	"ElmoBeacon/service"
	"cmp"
	"slices"
)

// GetItemHistory preserves sequence counts and annotations, and links a preceding
// loss with an explicit attribution status. Unfinished pity is excluded.
func (a *App) GetItemHistory(userID, itemID int64) ([]DisplayRecord, error) {
	var stored []model.Record
	if err := db.Engine.Distinct("pool_type").Where("user_id = ? AND item_id = ?", userID, itemID).Find(&stored); err != nil {
		return nil, err
	}
	types, err := a.GetPoolTypes(userID)
	if err != nil {
		return nil, err
	}
	result := []DisplayRecord{}
	for _, category := range types {
		if category == 9 || category == 10 {
			continue
		}
		if !slices.ContainsFunc(stored, func(record model.Record) bool { return record.PoolType == service.APIPoolType(category) }) {
			continue
		}
		info, err := a.GetPoolInfo(userID, category)
		if err != nil {
			return nil, err
		}
		chronological := slices.Clone(info.RecordList)
		slices.SortFunc(chronological, compareHistoryRecords)
		for _, record := range withAcquisitionCosts(chronological) {
			if record.ItemID == itemID {
				result = append(result, record)
			}
		}
	}
	slices.SortFunc(result, compareHistoryRecords)
	return result, nil
}

func compareHistoryRecords(a, b DisplayRecord) int {
	if order := cmp.Compare(a.Timestamp, b.Timestamp); order != 0 {
		return order
	}
	if order := cmp.Compare(len(a.ID), len(b.ID)); order != 0 {
		return order
	}
	return cmp.Compare(a.ID, b.ID)
}
