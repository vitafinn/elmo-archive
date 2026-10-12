package service

import (
	"ElmoBeacon/util"
	"fmt"
	"google.golang.org/protobuf/encoding/protowire"
	"strconv"
	"strings"
)

type GachaPoolInfo struct {
	PoolType  int64
	NameID    int64
	StartTime int64
	EndTime   int64
	UpItem    map[int64]struct{}
}

// rankItems reads IDs only, ignoring rate decimals and malformed tokens.
func rankItems(data string) map[int64]int {
	result := make(map[int64]int)
	for _, group := range strings.Split(strings.ReplaceAll(data, "：", ":"), ",") {
		parts := strings.Split(strings.TrimSpace(group), ":")
		rank, err := strconv.Atoi(parts[0])
		if err != nil || rank < 3 || rank > 5 {
			continue
		}
		for _, token := range parts[1:] {
			id, err := strconv.ParseInt(strings.TrimSpace(token), 10, 64)
			if err == nil && id >= 1000 {
				result[id] = rank
			}
		}
	}
	return result
}

// Read only fields needed for statistics. In newer clients field 35 became a
// string, whereas the legacy generated message treats it as packed item IDs.
func gachaFields(data []byte) (map[protowire.Number]int64, map[protowire.Number]string, error) {
	numbers := make(map[protowire.Number]int64)
	texts := make(map[protowire.Number]string)
	for len(data) > 0 {
		number, kind, n := protowire.ConsumeTag(data)
		if n < 0 {
			return nil, nil, protowire.ParseError(n)
		}
		data = data[n:]
		size := protowire.ConsumeFieldValue(number, kind, data)
		if size < 0 {
			return nil, nil, protowire.ParseError(size)
		}
		switch kind {
		case protowire.VarintType:
			value, _ := protowire.ConsumeVarint(data)
			numbers[number] = int64(value)
		case protowire.BytesType:
			value, _ := protowire.ConsumeBytes(data)
			texts[number] = string(value)
		}
		data = data[size:]
	}
	return numbers, texts, nil
}

func parseGachaPools(data []byte) (map[int64]GachaPoolInfo, map[int64]int, error) {
	pools := make(map[int64]GachaPoolInfo)
	ranks := make(map[int64]int)
	for len(data) > 0 {
		field, kind, n := protowire.ConsumeTag(data)
		if n < 0 {
			return nil, nil, protowire.ParseError(n)
		}
		data = data[n:]
		size := protowire.ConsumeFieldValue(field, kind, data)
		if size < 0 {
			return nil, nil, protowire.ParseError(size)
		}
		if field == 1 && kind == protowire.BytesType {
			row, _ := protowire.ConsumeBytes(data)
			numbers, texts, err := gachaFields(row)
			if err != nil {
				return nil, nil, err
			}
			id, typ := numbers[1], numbers[4]
			if id <= 0 || typ <= 0 {
				return nil, nil, fmt.Errorf("gacha row has no valid ID or type")
			}
			info := GachaPoolInfo{PoolType: typ, StartTime: numbers[7], EndTime: numbers[8], UpItem: make(map[int64]struct{})}
			if name := texts[9]; name != "" {
				nameFields, _, err := gachaFields([]byte(name))
				if err != nil {
					return nil, nil, err
				}
				info.NameID = nameFields[1]
			}
			// Current layout: all items 22/23, featured dolls 27, weapons 29.
			// Legacy layout: all items 21/22, featured dolls 29, weapons 28.
			allFields, upFields := []protowire.Number{22, 23}, []protowire.Number{27, 29}
			if len(rankItems(texts[21])) > 0 || len(rankItems(texts[28])) > 0 {
				allFields, upFields = []protowire.Number{21, 22}, []protowire.Number{29, 28}
			}
			for _, f := range allFields {
				for itemID, rank := range rankItems(texts[f]) {
					ranks[itemID] = rank
				}
			}
			// Selector tables contain a default choice, not the player's historical choice.
			// Preserve the existing standard-item fallback for selector pools 6 and 7.
			if typ != 6 && typ != 7 {
				for _, f := range upFields {
					for itemID, rank := range rankItems(texts[f]) {
						if rank == 5 {
							info.UpItem[itemID] = struct{}{}
						}
					}
				}
			}
			pools[id] = info
		}
		data = data[size:]
	}
	return pools, ranks, nil
}

func GetGachaPoolMap(gameDataDir, gameServer string) (map[int64]GachaPoolInfo, map[int64]int, error) {
	if gameServer == string(GameServerCN) {
		gameServer = ""
	}
	data, err := util.ReadTableData(gameDataDir, gameServer, "GachaData")
	if err != nil {
		return nil, nil, err
	}
	return parseGachaPools(data)
}
