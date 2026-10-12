package service

import (
	"ElmoBeacon/db"
	"ElmoBeacon/model"
	"ElmoBeacon/request"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"slices"
	"strconv"
)

type SyncDiff struct {
	PoolType int64
	Count    int
}

type SyncResult struct {
	Server   string
	Uid      uint64
	DiffList []SyncDiff
}

type ImportProgress struct {
	Stage    string
	PoolType int64
	Page     int
	Count    int
}

func SyncRecords(gameUserInfo *GameUserInfo, callbacks ...func(ImportProgress)) (*SyncResult, error) {
	if gameUserInfo.GameServer == GameServerCN {
		return nil, errors.New("the China client is not supported by Elmo Archive")
	}
	progress := func(stage string, pool int64, page, count int) {
		for _, callback := range callbacks {
			callback(ImportProgress{stage, pool, page, count})
		}
	}
	progress("syncStarted", 0, 0, 0)
	syncResult := SyncResult{
		Server: string(gameUserInfo.GameServer),
		Uid:    gameUserInfo.Uid,
	}

	// check if the user exists
	cond := model.User{GameServer: string(gameUserInfo.GameServer), Uid: gameUserInfo.Uid} //double conditions prevent users from having the same UID on different servers.It seems highly improbable, but it does exist among the users registered at the start of different servers.
	hasUser, err := db.Engine.Get(&cond)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, errors.New("Failed to query user")
	}
	var userId int64
	user := model.User{
		Uid:         gameUserInfo.Uid,
		GameServer:  string(gameUserInfo.GameServer),
		GameDataDir: gameUserInfo.GameDataDir,
	}
	if !hasUser {
		_, err = db.Engine.Insert(&user)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, errors.New("Failed to insert user")
		}
		userId = user.Id
	} else {
		userId = cond.Id
		_, err = db.Engine.ID(userId).Update(user)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, errors.New("Failed to update user")
		}
	}

	// fetch gacha records from official server until it matches the latest local record
	gachaPoolTypeList, err := GetGachaPoolTypeList(gameUserInfo.GameDataDir, string(gameUserInfo.GameServer))
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, errors.New("Failed to get gacha pool type list")
	}

	poolMap, _, metadataErr := GetGachaPoolMap(gameUserInfo.GameDataDir, string(gameUserInfo.GameServer))
	if metadataErr != nil {
		log.Warn().Err(metadataErr).Msg("Using known banner IDs for sync summary")
	}
	for _, poolType := range gachaPoolTypeList {
		progress("poolStarted", poolType, 0, 0)
		latestLocalRecord := model.Record{
			UserId:   userId,
			PoolType: poolType,
		}
		_, err = db.Engine.Desc("id").Get(&latestLocalRecord)
		if err != nil {
			log.Error().Err(err).Msg("")
			return nil, errors.New("Failed to query latest local record")
		}

		var incrementalRecordList []model.Record
		var next string
		page := 0
	loopFetchingRemoteRecord:
		for {
			page++
			remoteRecordList, err := request.FetchGachaRecordList(gameUserInfo.GachaRecordUrl, gameUserInfo.GameAccessToken, next, poolType)
			if err != nil {
				log.Error().Err(err).Msg("")
				return nil, errors.New("Failed to fetch gacha record list")
			}
			progress("pageFetched", poolType, page, len(remoteRecordList.RecordList))
			for _, remoteRecord := range remoteRecordList.RecordList {
				if remoteRecord.PoolId == latestLocalRecord.PoolId && remoteRecord.GachaTimestamp == latestLocalRecord.Timestamp && remoteRecord.ItemId == latestLocalRecord.ItemId {
					break loopFetchingRemoteRecord
				} else {
					incrementalRecordList = append(incrementalRecordList, model.Record{
						UserId:    userId,
						PoolType:  poolType,
						PoolId:    remoteRecord.PoolId,
						ItemId:    remoteRecord.ItemId,
						Timestamp: remoteRecord.GachaTimestamp,
					})
				}
			}
			if remoteRecordList.Next != "" {
				next = remoteRecordList.Next
			} else {
				break
			}
		}
		progress("poolFetched", poolType, page, len(incrementalRecordList))

		// merge the incremental gacha records into the database
		if len(incrementalRecordList) > 0 {
			slices.Reverse(incrementalRecordList)
			var lastTimestamp, order int64
			for i, record := range incrementalRecordList {
				if record.Timestamp != lastTimestamp {
					order = 0
				} else {
					order++
				}
				virtualId, _ := strconv.ParseUint(fmt.Sprintf("%d%03d", record.Timestamp, order), 10, 64)
				incrementalRecordList[i].Id = virtualId
				lastTimestamp = record.Timestamp
			}

			_, err = db.Engine.Insert(&incrementalRecordList)
			if err != nil {
				log.Error().Err(err).Msg("")
				return nil, errors.New("Failed to insert incremental record list")
			}

			log.Info().Str("server", string(gameUserInfo.GameServer)).Uint64("uid", gameUserInfo.Uid).Int64("poolType", poolType).Int("count", len(incrementalRecordList)).Msg("")

			counts := make(map[int64]int)
			for _, record := range incrementalRecordList {
				counts[DisplayPoolType(record.PoolType, record.PoolId, poolMap)]++
			}
			categories := make([]int64, 0, len(counts))
			for category := range counts {
				categories = append(categories, category)
			}
			slices.Sort(categories)
			for _, category := range categories {
				syncResult.DiffList = append(syncResult.DiffList, SyncDiff{PoolType: category, Count: counts[category]})
			}
		}
	}
	progress("syncDone", 0, 0, 0)

	return &syncResult, nil
}
