package service

// Negative IDs are display-only categories, never remote API type IDs.
const (
	ReunionDolls   int64 = -6
	ReunionWeapons int64 = -7
)

func APIPoolType(displayType int64) int64 {
	switch displayType {
	case ReunionDolls:
		return 6
	case ReunionWeapons:
		return 7
	default:
		return displayType
	}
}

// Reunion and Custom share API types 6/7. The client gives Reunion its own
// localized name IDs. Known banner IDs keep historical records separated when
// a later client drops those banners from its tables.
func DisplayPoolType(apiType, bannerID int64, pools map[int64]GachaPoolInfo) int64 {
	info, known := pools[bannerID]
	if apiType == 6 && ((known && info.NameID == 559138) || bannerID == 252001) {
		return ReunionDolls
	}
	if apiType == 7 && ((known && info.NameID == 559139) || bannerID == 253001) {
		return ReunionWeapons
	}
	return apiType
}

func DisplayPoolTypes(apiTypes []int64, pools map[int64]GachaPoolInfo) []int64 {
	result := NormalizePoolTypes(apiTypes)
	seen := make(map[int64]bool)
	for bannerID, info := range pools {
		seen[DisplayPoolType(info.PoolType, bannerID, pools)] = true
	}
	for _, typ := range []int64{ReunionDolls, ReunionWeapons} {
		if seen[typ] {
			result = append(result, typ)
		}
	}
	return result
}
