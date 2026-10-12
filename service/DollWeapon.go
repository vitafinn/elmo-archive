package service

import (
	"ElmoBeacon/pb"
	"ElmoBeacon/util"
	"fmt"
	"google.golang.org/protobuf/encoding/protowire"
)

// GunData field 24 is the elite weapon link, not the general recommended
// weapon list (which also recommends shared and battle-pass weapons).
func GetDollWeaponMap(dir, server string) (map[int64]int64, error) {
	data, err := util.ReadTableData(dir, server, "GunData")
	if err != nil {
		return nil, err
	}
	links, err := parseDollWeapons(data)
	if err != nil {
		return nil, err
	}
	var weapons pb.GunWeaponData
	if err = util.GetTableData(dir, server, &weapons); err != nil {
		return nil, err
	}
	items, err := GetItemDataMap(dir, server)
	if err != nil {
		return nil, err
	}
	known := map[int64]bool{}
	for _, weapon := range weapons.Units {
		known[weapon.Id] = true
	}
	for doll, weapon := range links {
		if items[doll].GetType() != 10 || !known[weapon] || items[weapon].GetType() != 20 {
			delete(links, doll)
		}
	}
	return links, nil
}

func parseDollWeapons(data []byte) (map[int64]int64, error) {
	result := map[int64]int64{}
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
			numbers, _, err := gachaFields(row)
			if err != nil {
				return nil, err
			}
			doll, weapon := numbers[1], numbers[24]
			if doll > 0 && weapon > 0 {
				if previous, exists := result[doll]; exists && previous != weapon {
					return nil, fmt.Errorf("conflicting weapon links for doll %d", doll)
				}
				result[doll] = weapon
			}
		}
		data = data[size:]
	}
	return result, nil
}
