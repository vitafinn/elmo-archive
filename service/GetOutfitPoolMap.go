package service

import (
	"ElmoBeacon/util"
	"google.golang.org/protobuf/encoding/protowire"
	"strconv"
	"strings"
)

type OutfitPoolInfo struct {
	MainItem int64
	Rewards  map[int64]string
}

func outfitRewards(text, category string, result map[int64]string) {
	for _, group := range strings.Split(strings.ReplaceAll(text, "：", ":"), ",") {
		parts := strings.Split(strings.TrimSpace(group), ":")
		if len(parts) < 2 {
			continue
		}
		rank, err := strconv.Atoi(parts[0])
		if err != nil || rank < 3 || rank > 5 {
			continue
		}
		for _, token := range parts[1:] {
			id, err := strconv.ParseInt(strings.TrimSpace(token), 10, 64)
			if err != nil || id <= 0 {
				continue
			}
			kind := category
			if kind == "" {
				if rank >= 4 {
					kind = "rare"
				} else {
					kind = "normal"
				}
			}
			result[id] = kind
		}
	}
}

func parseOutfitPools(data []byte) (map[int64]OutfitPoolInfo, error) {
	result := make(map[int64]OutfitPoolInfo)
	for len(data) > 0 {
		field, kind, n := protowire.ConsumeTag(data)
		if n < 0 {
			return nil, protowire.ParseError(n)
		}
		data = data[n:]
		size := protowire.ConsumeFieldValue(field, kind, data)
		if size < 0 {
			return nil, protowire.ParseError(size)
		}
		if field == 1 && kind == protowire.BytesType {
			row, _ := protowire.ConsumeBytes(data)
			numbers, texts, err := gachaFields(row)
			if err != nil {
				return nil, err
			}
			if numbers[1] > 0 {
				info := OutfitPoolInfo{MainItem: numbers[11], Rewards: make(map[int64]string)}
				outfitRewards(texts[9], "", info.Rewards)
				outfitRewards(texts[7], "limited", info.Rewards)
				result[numbers[1]] = info
			}
		}
		data = data[size:]
	}
	return result, nil
}

func GetOutfitPoolMap(dir, server string) (map[int64]OutfitPoolInfo, error) {
	if server == string(GameServerCN) {
		server = ""
	}
	data, err := util.ReadTableData(dir, server, "GachaClothesData")
	if err != nil {
		return nil, err
	}
	return parseOutfitPools(data)
}
