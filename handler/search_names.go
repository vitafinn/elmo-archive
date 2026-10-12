package handler

import (
	"ElmoBeacon/pb"
	"ElmoBeacon/service"
	"slices"
	"strings"
)

// Missing optional language files must not prevent history from loading.
func searchLanguages(dir, server string, current map[int64]string) []map[int64]string {
	languages := []map[int64]string{current}
	for _, locale := range []string{"en", "zh-TW"} {
		if names, err := service.GetLangDataMap(dir, server, locale); err == nil {
			languages = append(languages, names)
		}
	}
	return languages
}

func itemSearchNames(id int64, items map[int64]*pb.Item, languages []map[int64]string) []string {
	names := make([]string, 0, len(languages))
	item := items[id]
	if item == nil || item.Name == nil {
		return names
	}
	for _, language := range languages {
		name := strings.TrimSpace(language[item.Name.Id])
		if name != "" && !slices.Contains(names, name) {
			names = append(names, name)
		}
	}
	return names
}
