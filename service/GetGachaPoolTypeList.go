package service

import (
	"ElmoBeacon/pb"
	"ElmoBeacon/util"
	"slices"
)

// NormalizePoolTypes preserves new categories without duplicating requests or tabs.
func NormalizePoolTypes(groups ...[]int64) []int64 {
	seen := make(map[int64]bool)
	result := make([]int64, 0)
	for _, group := range groups {
		for _, id := range group {
			if id > 0 && !seen[id] {
				seen[id] = true
				result = append(result, id)
			}
		}
	}
	slices.Sort(result)
	return result
}

func GetGachaPoolTypeList(gameDataDir, gameServer string) ([]int64, error) {
	var data pb.GachaTypeListData
	if gameServer == string(GameServerCN) {
		gameServer = ""
	}
	if err := util.GetTableData(gameDataDir, gameServer, &data); err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(data.Units))
	for _, unit := range data.Units {
		if unit != nil {
			ids = append(ids, unit.Id)
		}
	}
	return NormalizePoolTypes(ids), nil
}
