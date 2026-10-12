package handler

import (
	"ElmoBeacon/db"
	"ElmoBeacon/model"
	"ElmoBeacon/service"
	"fmt"
)

type DollWeaponHistory struct {
	ItemID  int64
	Name    string
	Records []DisplayRecord
}

// A missing link or missing history is distinct from a weapon that cost zero pulls.
func (a *App) GetDollWeaponHistory(userID, dollID int64) (DollWeaponHistory, error) {
	result := DollWeaponHistory{Records: []DisplayRecord{}}
	var user model.User
	found, err := db.Engine.ID(userID).Get(&user)
	if err != nil {
		return result, err
	}
	if !found {
		return result, fmt.Errorf("account not found")
	}
	links, err := service.GetDollWeaponMap(user.GameDataDir, user.GameServer)
	if err != nil {
		return result, err
	}
	result.ItemID = links[dollID]
	if result.ItemID == 0 {
		return result, nil
	}
	lang := model.Setting{Key: "lang"}
	if _, err = db.Engine.Get(&lang); err != nil {
		return result, err
	}
	items, err := service.GetItemDataMap(user.GameDataDir, user.GameServer)
	if err != nil {
		return result, err
	}
	names, err := service.GetLangDataMap(user.GameDataDir, user.GameServer, lang.Value)
	if err != nil {
		return result, err
	}
	result.Name = names[items[result.ItemID].GetName().GetId()]
	if result.Name == "" {
		result.Name = fmt.Sprintf("Weapon #%d", result.ItemID)
	}
	result.Records, err = a.GetItemHistory(userID, result.ItemID)
	return result, err
}
