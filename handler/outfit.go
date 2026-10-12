package handler

import (
	"ElmoBeacon/db"
	"ElmoBeacon/model"
	"ElmoBeacon/pb"
	"ElmoBeacon/service"
	"fmt"
	"github.com/rs/zerolog/log"
	"slices"
	"strconv"
	"time"
)

type OutfitRecord struct {
	ID           string   `json:"id"`
	ItemID       int64    `json:"itemId"`
	Name         string   `json:"name"`
	SearchNames  []string `json:"searchNames"`
	Category     string   `json:"category"`
	AccessNumber int64    `json:"accessNumber"`
	Timestamp    int64    `json:"timestamp"`
}
type OutfitSeason struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	StartTime    int64          `json:"startTime"`
	EndTime      int64          `json:"endTime"`
	Status       string         `json:"status"`
	TotalCount   int64          `json:"totalCount"`
	LimitedCount int64          `json:"limitedCount"`
	RareCount    int64          `json:"rareCount"`
	NormalCount  int64          `json:"normalCount"`
	UnknownCount int64          `json:"unknownCount"`
	MainOutfit   string         `json:"mainOutfit"`
	MainOutfitAt int64          `json:"mainOutfitAt"`
	Records      []OutfitRecord `json:"records"`
}
type OutfitInfo struct {
	TotalCount int64          `json:"totalCount"`
	Seasons    []OutfitSeason `json:"seasons"`
}

func outfitItemName(id int64, items map[int64]*pb.Item, lang map[int64]string) string {
	if item := items[id]; item != nil && item.Name != nil {
		if name := lang[item.Name.Id]; name != "" {
			return name
		}
	}
	return fmt.Sprintf("Item #%d", id)
}

// Every access is retained, including currency and rewards absent from client tables.
func buildOutfitInfo(records []model.Record, poolType int64, pools map[int64]service.GachaPoolInfo, outfits map[int64]service.OutfitPoolInfo, items map[int64]*pb.Item, lang map[int64]string, now int64) OutfitInfo {
	info := OutfitInfo{Seasons: make([]OutfitSeason, 0)}
	seasons := make(map[int64]*OutfitSeason)
	ensure := func(id int64) *OutfitSeason {
		if existing := seasons[id]; existing != nil {
			return existing
		}
		meta := pools[id]
		season := &OutfitSeason{ID: id, Name: lang[meta.NameID], StartTime: meta.StartTime, EndTime: meta.EndTime, Status: "unknown", Records: make([]OutfitRecord, 0)}
		if season.Name == "" {
			season.Name = fmt.Sprintf("Outfit Procurement #%d", id)
		}
		if meta.EndTime > 0 && now > meta.EndTime {
			season.Status = "expired"
		} else if meta.StartTime > 0 && meta.StartTime <= now && now <= meta.EndTime {
			season.Status = "active"
		}
		if outfit := outfits[id]; outfit.MainItem > 0 {
			season.MainOutfit = outfitItemName(outfit.MainItem, items, lang)
		}
		seasons[id] = season
		return season
	}
	for id, meta := range pools {
		if meta.PoolType == poolType && meta.StartTime > 0 && meta.StartTime <= now && now <= meta.EndTime {
			ensure(id)
		}
	}
	for _, record := range records {
		if record.PoolType != poolType {
			continue
		}
		season := ensure(record.PoolId)
		season.TotalCount++
		info.TotalCount++
		category := outfits[record.PoolId].Rewards[record.ItemId]
		switch category {
		case "limited":
			season.LimitedCount++
		case "rare":
			season.RareCount++
		case "normal":
			season.NormalCount++
		default:
			category = "unknown"
			season.UnknownCount++
		}
		if record.ItemId == outfits[record.PoolId].MainItem && season.MainOutfitAt == 0 {
			season.MainOutfitAt = season.TotalCount
		}
		season.Records = append(season.Records, OutfitRecord{ID: strconv.FormatUint(record.Id, 10), ItemID: record.ItemId, Name: outfitItemName(record.ItemId, items, lang), Category: category, AccessNumber: season.TotalCount, Timestamp: record.Timestamp})
	}
	for _, season := range seasons {
		slices.Reverse(season.Records)
		info.Seasons = append(info.Seasons, *season)
	}
	slices.SortFunc(info.Seasons, func(a, b OutfitSeason) int {
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
	return info
}
func (a *App) GetOutfitInfo(userID, poolType int64) (OutfitInfo, error) {
	if poolType != 9 && poolType != 10 {
		return OutfitInfo{}, fmt.Errorf("not an outfit category")
	}
	var user model.User
	found, err := db.Engine.ID(userID).Get(&user)
	if err != nil {
		return OutfitInfo{}, err
	}
	if !found {
		return OutfitInfo{}, fmt.Errorf("user not found")
	}
	var records []model.Record
	if err := db.Engine.Where("user_id = ? AND pool_type = ?", userID, poolType).OrderBy("id").Find(&records); err != nil {
		return OutfitInfo{}, err
	}
	langSetting, err := a.GetSetting("lang")
	if err != nil {
		return OutfitInfo{}, err
	}
	lang, err := service.GetLangDataMap(user.GameDataDir, user.GameServer, langSetting)
	if err != nil {
		log.Warn().Err(err).Msg("Outfit labels unavailable; displaying IDs")
	}
	items, err := service.GetItemDataMap(user.GameDataDir, user.GameServer)
	if err != nil {
		log.Warn().Err(err).Msg("Outfit item data unavailable; displaying IDs")
	}
	pools, _, err := service.GetGachaPoolMap(user.GameDataDir, user.GameServer)
	if err != nil {
		log.Warn().Err(err).Msg("Outfit schedule unavailable")
	}
	outfits, err := service.GetOutfitPoolMap(user.GameDataDir, user.GameServer)
	if err != nil {
		log.Warn().Err(err).Msg("Outfit reward categories unavailable")
	}
	info := buildOutfitInfo(records, poolType, pools, outfits, items, lang, time.Now().Unix())
	languages := searchLanguages(user.GameDataDir, user.GameServer, lang)
	for i := range info.Seasons {
		for j := range info.Seasons[i].Records {
			record := &info.Seasons[i].Records[j]
			record.SearchNames = itemSearchNames(record.ItemID, items, languages)
		}
	}
	return info, nil
}
