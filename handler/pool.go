package handler

import (
	"ElmoBeacon/db"
	"ElmoBeacon/model"
	"ElmoBeacon/service"
	"slices"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type DisplayRecord struct {
	ID               string
	ItemID           int64
	ItemType         int64
	PoolType         int64
	BannerID         int64
	BannerName       string
	Outcome          string
	Annotation       string
	Name             string
	SearchNames      []string
	Icon             string
	Count            int64
	Timestamp        int64
	IsMissing        bool
	TargetKnown      bool
	AcquisitionCount int64
	PreviousLoss     *AcquisitionLoss
	CostBasis        string
}

type PoolInfo struct {
	BannerPeriods    []BannerPeriod  `json:"bannerPeriods"`
	Rank5Spent       float64         `json:"rank5Spent"`
	ExpiredPity      int64           `json:"expiredPity"`
	PendingPity      int64           `json:"pendingPity"`
	UnknownRankCount int64           `json:"unknownRankCount"`
	TargetAttempts   int64           `json:"targetAttempts"`
	TargetWins       int64           `json:"targetWins"`
	UnknownOutcomes  int64           `json:"unknownOutcomes"`
	StoredCount      int64           `json:"storedCount"`
	RecordList       []DisplayRecord `json:"recordList"`
	TotalCount       int64           `json:"totalCount"`
	Rank5Count       int64           `json:"rank5Count"`
	Rank4Count       int64           `json:"rank4Count"`
	Rank3Count       int64           `json:"rank3Count"`
	Rank5Rate        float64         `json:"rank5Rate"`
	Rank4Rate        float64         `json:"rank4Rate"`
	Rank3Rate        float64         `json:"rank3Rate"`
	Rank5Avg         float64         `json:"rank5Avg"`
	Rank5UpAvg       float64         `json:"rank5UpAvg"`
	MissingCount     int64           `json:"missingCount"`
	MissingRate      float64         `json:"missingRate"`
}

func getItemRank(id int64) int {
	fourStarDolls := map[int64]bool{
		1001: true, 1008: true, 1009: true, 1017: true, 1022: true,
		1024: true, 1026: true, 1036: true, 1038: true, 1041: true,
	}
	fiveStarWeapons11 := map[int64]bool{
		11016: true, 11020: true, 11038: true, 11044: true, 11047: true, 11053: true, 11056: true,
	}
	fourStarWeapons11 := map[int64]bool{
		11007: true, 11014: true, 11015: true, 11021: true, 11023: true, 11026: true, 11031: true,
		11037: true, 11040: true, 11043: true, 11046: true, 11049: true, 11052: true, 11055: true,
	}

	// Strict boundaries applied. Ignored items (like 9010) drop to bottom fallback.
	if id >= 1000 && id < 2000 {
		if fourStarDolls[id] {
			return 4
		}
		return 5
	}

	if id >= 10000 && id < 20000 {
		if id >= 10001 && id <= 10099 {
			return 5
		}
		if id >= 10100 && id < 11000 {
			lastDigit := id % 10
			if lastDigit == 1 {
				return 3
			}
			if lastDigit == 2 {
				return 4
			}
			if lastDigit == 3 {
				return 5
			}
		}
		if id >= 11000 {
			if fiveStarWeapons11[id] {
				return 5
			}
			if fourStarWeapons11[id] {
				return 4
			}
			return 3
		}
	}

	return 0
}

func isStandard5Star(id int64) bool {
	std := map[int64]bool{
		1015: true, 1021: true, 1025: true, 1027: true, 1028: true, 1029: true, 1033: true,
		11016: true, 11020: true, 11038: true, 11044: true, 11047: true, 11053: true, 11056: true, 10333: true,
	}
	return std[id]
}

func (a *App) GetPoolInfo(userId, poolType int64) (poolInfo PoolInfo, err error) {
	var user model.User
	_, err = db.Engine.ID(userId).Get(&user)
	if err != nil {
		return PoolInfo{}, err
	}
	lang := model.Setting{Key: "lang"}
	_, err = db.Engine.Get(&lang)
	if err != nil {
		return PoolInfo{}, err
	}

	langDataMap, err := service.GetLangDataMap(user.GameDataDir, user.GameServer, lang.Value)
	if err != nil {
		log.Error().Err(err).Msg("")
		return PoolInfo{}, errors.New("error occurred when get lang map")
	}

	itemDataMap, err := service.GetItemDataMap(user.GameDataDir, user.GameServer)
	if err != nil {
		log.Error().Err(err).Msg("")
		return PoolInfo{}, errors.New("error occurred when get item map")
	}

	searchNames := searchLanguages(user.GameDataDir, user.GameServer, langDataMap)
	dollIconMap, err := service.GetDollIconMap(user.GameDataDir, user.GameServer)
	if err != nil {
		log.Error().Err(err).Msg("")
		return PoolInfo{}, errors.New("error occurred when get doll icon map")
	}

	weaponIconMap, err := service.GetWeaponIconMap(user.GameDataDir, user.GameServer)
	if err != nil {
		log.Error().Err(err).Msg("")
		return PoolInfo{}, errors.New("error occurred when get weapon icon map")
	}

	gachaPoolMap, tableRanks, err := service.GetGachaPoolMap(user.GameDataDir, user.GameServer)
	if err != nil {
		log.Error().Err(err).Msg("")
		return PoolInfo{}, errors.New("error occurred when get gacha pool map")
	}

	var recordList []model.Record
	err = db.Engine.OrderBy("id").Where("user_id = ? AND pool_type = ?", userId, service.APIPoolType(poolType)).Find(&recordList)
	if err != nil {
		log.Error().Err(err).Msg("")
		return PoolInfo{}, errors.New("error occurred when get user record list by poolType from db")
	}

	isSelector := service.APIPoolType(poolType) == 6 || service.APIPoolType(poolType) == 7
	annotations := make(map[string]string)
	if isSelector {
		annotations, err = a.selectorAnnotations(userId)
		if err != nil {
			return PoolInfo{}, err
		}
	}
	periods := make(map[int64]*BannerPeriod)
	if isSelector {
		for id, metadata := range gachaPoolMap {
			if service.DisplayPoolType(metadata.PoolType, id, gachaPoolMap) == poolType && metadata.StartTime <= time.Now().Unix() && metadata.EndTime >= time.Now().Unix() {
				periods[id] = newBannerPeriod(id, metadata)
			}
		}
	}
	if len(recordList) > 0 {
		var displayRecordList []DisplayRecord
		var count, judgeCount int64
		var isPreMissing bool
		selectorHistory := make(selectorGuarantees)
		selectorCounts := make(map[int64]int64)
		var completedPulls int64

		for _, record := range recordList {
			if service.DisplayPoolType(record.PoolType, record.PoolId, gachaPoolMap) != poolType {
				continue
			}
			item, hasItemData := itemDataMap[record.ItemId]
			if !hasItemData && !isSelector {
				continue
			}
			rank := 0
			if hasItemData {
				rank = int(item.Rank)
			}
			if rank == 0 {
				rank = tableRanks[record.ItemId]
			}
			if rank == 0 {
				rank = getItemRank(record.ItemId)
			}
			if rank == 0 && !isSelector {
				rank = 3
			}

			var period *BannerPeriod
			if isSelector {
				period = periods[record.PoolId]
				if period == nil {
					period = newBannerPeriod(record.PoolId, gachaPoolMap[record.PoolId])
					periods[record.PoolId] = period
				}
				period.TotalCount++
				period.LastTimestamp = record.Timestamp
			}
			poolInfo.TotalCount++
			if isSelector {
				count = selectorCounts[record.PoolId]
			}
			count++

			if rank == 5 {
				text := "Unknown #" + strconv.FormatInt(record.ItemId, 10)
				if hasItemData && item.Name != nil {
					if langText, hasLangData := langDataMap[item.Name.Id]; hasLangData {
						text = langText
					}
				}

				poolInfo.Rank5Count++

				var isMissing, targetKnown bool
				outcome, annotation := "", ""
				if isSelector {
					annotation = annotations[strconv.FormatUint(record.Id, 10)]
					var eligible bool
					outcome, eligible = selectorHistory.observe(record.PoolId, annotation)
					isMissing = outcome == "loss"
					if eligible {
						poolInfo.TargetAttempts++
						period.TargetAttempts++
					}
					if outcome == "win" {
						poolInfo.TargetWins++
						period.TargetWins++
					}
					if outcome == "unknown" {
						poolInfo.UnknownOutcomes++
						period.UnknownOutcomes++
					}
					if isMissing {
						poolInfo.MissingCount++
					}
					period.Rank5Count++
					period.CompletedPulls += count
					completedPulls += count
				} else if record.PoolType == 3 || record.PoolType == 4 {
					if !isPreMissing {
						judgeCount++
					}
					info, known := gachaPoolMap[record.PoolId]
					if known && len(info.UpItem) > 0 {
						targetKnown = true
						_, isUp := info.UpItem[item.Id]
						isMissing = !isUp
					} else {
						isMissing = isStandard5Star(item.Id)
					}
					if isMissing {
						poolInfo.MissingCount++
					}
				}

				icon := "Avatar_Head_Unknown_Spine.png"
				if hasItemData && item.Type == 10 {
					if i, hasIcon := dollIconMap[item.Id]; hasIcon {
						icon = i
					}
				} else if hasItemData && item.Type == 20 {
					if i, hasIcon := weaponIconMap[item.Id]; hasIcon {
						icon = i
					}
				}

				displayRecordList = append(displayRecordList, DisplayRecord{
					ID:          strconv.FormatUint(record.Id, 10),
					ItemID:      record.ItemId,
					ItemType:    int64(item.GetType()),
					PoolType:    poolType,
					BannerID:    record.PoolId,
					BannerName:  langDataMap[gachaPoolMap[record.PoolId].NameID],
					Outcome:     outcome,
					Annotation:  annotation,
					Name:        text,
					SearchNames: itemSearchNames(record.ItemId, itemDataMap, searchNames),
					Icon:        icon,
					Count:       count,
					Timestamp:   record.Timestamp,
					IsMissing:   isMissing,
					TargetKnown: targetKnown,
				})

				isPreMissing = isMissing
				count = 0
			} else if rank == 4 {
				poolInfo.Rank4Count++
			} else if rank == 3 {
				poolInfo.Rank3Count++
			}
			if isSelector {
				selectorCounts[record.PoolId] = count
				period.TailCount = count
				if rank == 0 {
					poolInfo.UnknownRankCount++
					period.UnknownRankCount++
				}
			}
		}

		poolInfo.StoredCount = count
		if judgeCount > 0 {
			poolInfo.MissingRate = float64(poolInfo.MissingCount) / float64(judgeCount)
		}
		if poolInfo.TotalCount > 0 {
			poolInfo.Rank5Rate = float64(poolInfo.Rank5Count) / float64(poolInfo.TotalCount)
			poolInfo.Rank4Rate = float64(poolInfo.Rank4Count) / float64(poolInfo.TotalCount)
			poolInfo.Rank3Rate = float64(poolInfo.Rank3Count) / float64(poolInfo.TotalCount)
		}

		if poolInfo.Rank5Count > 0 {
			poolInfo.Rank5Spent = float64(poolInfo.TotalCount) / float64(poolInfo.Rank5Count)
			poolInfo.Rank5Avg = poolInfo.Rank5Spent
			if isSelector {
				poolInfo.Rank5Avg = float64(completedPulls) / float64(poolInfo.Rank5Count)
			}
		}
		if poolInfo.Rank5Count-poolInfo.MissingCount > 0 && (!isSelector || poolInfo.UnknownOutcomes == 0) {
			poolInfo.Rank5UpAvg = float64(poolInfo.TotalCount) / float64(poolInfo.Rank5Count-poolInfo.MissingCount)
		}

		slices.Reverse(displayRecordList)
		poolInfo.RecordList = displayRecordList
	}

	if isSelector {
		finalizeSelectorPeriods(&poolInfo, periods, time.Now().Unix())
		for i := range poolInfo.BannerPeriods {
			period := &poolInfo.BannerPeriods[i]
			period.Name = langDataMap[gachaPoolMap[period.ID].NameID]
		}
	}
	return poolInfo, nil
}
