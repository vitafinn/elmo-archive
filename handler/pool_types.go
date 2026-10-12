package handler

import (
	"ElmoBeacon/db"
	"ElmoBeacon/model"
	"ElmoBeacon/service"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// GetPoolTypes includes historical categories even after the client removes them.
func (a *App) GetPoolTypes(userID int64) ([]int64, error) {
	fallback := []int64{1, 3, 4, 5, 6, 7, 8, 9}
	if userID == 0 {
		return append(fallback, service.ReunionDolls, service.ReunionWeapons), nil
	}
	var user model.User
	found, err := db.Engine.ID(userID).Get(&user)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("user not found")
	}
	var records []model.Record
	if err := db.Engine.Distinct("pool_type", "pool_id").Where("user_id = ?", userID).Find(&records); err != nil {
		return nil, err
	}
	stored := make([]int64, 0, len(records))
	for _, record := range records {
		stored = append(stored, record.PoolType)
	}
	discovered, err := service.GetGachaPoolTypeList(user.GameDataDir, user.GameServer)
	if err != nil {
		log.Warn().Err(err).Msg("Pool discovery unavailable; showing known and stored categories")
		discovered = fallback
	}
	pools, _, metadataErr := service.GetGachaPoolMap(user.GameDataDir, user.GameServer)
	if metadataErr != nil {
		pools = make(map[int64]service.GachaPoolInfo)
	}
	for _, record := range records {
		if _, exists := pools[record.PoolId]; !exists {
			pools[record.PoolId] = service.GachaPoolInfo{PoolType: record.PoolType}
		}
	}
	return service.DisplayPoolTypes(service.NormalizePoolTypes(discovered, stored), pools), nil
}
