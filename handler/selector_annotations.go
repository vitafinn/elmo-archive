package handler

import (
	"ElmoBeacon/db"
	"ElmoBeacon/model"
	"fmt"
	"strconv"
	"strings"
)

func selectorPrefix(userID int64) string { return fmt.Sprintf("selector.outcome/%d/", userID) }
func (a *App) selectorAnnotations(userID int64) (map[string]string, error) {
	prefix := selectorPrefix(userID)
	var settings []model.Setting
	if err := db.Engine.Where("key LIKE ?", prefix+"%").Find(&settings); err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, setting := range settings {
		result[strings.TrimPrefix(setting.Key, prefix)] = setting.Value
	}
	return result, nil
}

// Annotations never modify the original imported pull or its pity counter.
func (a *App) SetSelectorOutcome(userID int64, recordID, outcome string) error {
	switch outcome {
	case "auto", "unknown", "initial", "guaranteed", "win", "loss":
	default:
		return fmt.Errorf("invalid target outcome")
	}
	id, err := strconv.ParseUint(recordID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid record ID")
	}
	var record model.Record
	found, err := db.Engine.Where("id = ? AND user_id = ?", id, userID).Get(&record)
	if err != nil {
		return err
	}
	if !found || (record.PoolType != 6 && record.PoolType != 7) {
		return fmt.Errorf("selector record not found for this user")
	}
	return a.SetSetting(selectorPrefix(userID)+strconv.FormatUint(id, 10), outcome)
}
