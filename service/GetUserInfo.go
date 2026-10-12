package service

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type UserProvidedConfig struct {
	IdentitySource     string `json:"IdentitySource"`
	FullGachaRecordUrl string `json:"FullGachaRecordUrl"`
	AccessToken        string `json:"AccessToken"`
	UidString          string `json:"UidString"`
}

const CONFIG_FILE_NAME = "gacha_session_info.json"

func loadUserConfig(configFilePath string) (*UserProvidedConfig, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Wrapf(err, "config file '%s' not found", configFilePath)
		}
		return nil, errors.Wrapf(err, "failed to read config file '%s'", configFilePath)
	}

	var config UserProvidedConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal config data from file '%s'", configFilePath)
	}

	if config.FullGachaRecordUrl == "" {
		return nil, errors.New("config file error: missing 'FullGachaRecordUrl'")
	}
	if config.AccessToken == "" {
		return nil, errors.New("config file error: missing 'AccessToken'")
	}
	if config.UidString == "" || strings.ToUpper(config.UidString) == "REPLACE_WITH_YOUR_UID" {
		return nil, errors.New("config file error: 'UidString' is placeholder or missing. Please edit '" + CONFIG_FILE_NAME + "'.")
	}
	return &config, nil
}

func GetUserInfoOS() (*GameUserInfo, error) {
	// --- 1. RAM INJECTION PROTOCOL ---
	sessionMutex.Lock()
	memConfig := ActiveSession
	sessionMutex.Unlock()

	var authConfig *UserProvidedConfig

	if memConfig.FullGachaRecordUrl != "" && memConfig.AccessToken != "" && memConfig.UidString != "" {
		authConfig = &UserProvidedConfig{
			IdentitySource:     memConfig.IdentitySource,
			FullGachaRecordUrl: memConfig.FullGachaRecordUrl,
			AccessToken:        memConfig.AccessToken,
			UidString:          memConfig.UidString,
		}
	} else {
		// --- 2. ORIGINAL JSON FALLBACK ---
		exePath, errExe := os.Executable()
		configFilePath := CONFIG_FILE_NAME
		if errExe != nil {

		} else {
			configFilePath = filepath.Join(filepath.Dir(exePath), CONFIG_FILE_NAME)
		}

		var errConfig error
		authConfig, errConfig = loadUserConfig(configFilePath)
		if errConfig != nil {
			// Silently fail to trigger the final Player.log fallback
		}
	}

	if authConfig != nil {
		uid, parseErr := strconv.ParseUint(authConfig.UidString, 10, 64)
		if parseErr != nil {
			return nil, errors.Wrapf(parseErr, "Failed to parse UidString '%s' from JSON config (OS)", authConfig.UidString)
		}

		gachaRecordUrlFromConfig := authConfig.FullGachaRecordUrl
		var server gameServer

		urlLowerConfig := strings.ToLower(gachaRecordUrlFromConfig)
		switch {
		case strings.Contains(urlLowerConfig, "gf2-gacha-record-us"):
			server = GameServerUS
		case strings.Contains(urlLowerConfig, "gf2-gacha-record-intl"):
			server = GameServerGlobal
		case strings.Contains(urlLowerConfig, "gf2-gacha-record-jp"):
			server = GameServerJP
		case strings.Contains(urlLowerConfig, "gf2-gacha-record-kr"):
			server = GameServerKR
		case strings.Contains(urlLowerConfig, "gf2-gacha-record-asia"):
			server = GameServerAsia
		default:
			return nil, errors.Errorf("Failed to determine a valid Oversea server from JSON config URL: %s", gachaRecordUrlFromConfig)
		}

		// --- 3. THE ORIGINAL GAME DIRECTORY EXTRACTION ---
		var gameDataDir string
		userHome, homeErr := os.UserHomeDir()
		if homeErr != nil {

		} else {
			logDataBytes, readLogErr := os.ReadFile(filepath.Join(userHome, gameLogPathOversea))
			if readLogErr != nil {

			} else {
				regexpGameDataDir, compErrGD := regexp.Compile(exprGameDataDir)
				if compErrGD != nil {
					return nil, errors.New("Internal error: Failed to compile exprGameDataDir for GameDataDir")
				} else {
					resultGameDataDir := regexpGameDataDir.FindSubmatch(logDataBytes)
					if resultGameDataDir == nil {

					} else {
						gameDataDir = filepath.Join(string(resultGameDataDir[1]), "LocalCache/Data")
					}
				}
			}
		}

		userInfo := &GameUserInfo{
			IdentitySource:  authConfig.IdentitySource,
			Uid:             uid,
			GameServer:      server,
			GameDataDir:     gameDataDir,
			GameAccessToken: authConfig.AccessToken,
			GachaRecordUrl:  authConfig.FullGachaRecordUrl,
		}
		return userInfo, nil
	}

	userHome, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return nil, errors.New("Failed to get user home dir")
	}

	logDataBytes, err := os.ReadFile(filepath.Join(userHome, gameLogPathOversea))
	if err != nil {
		return nil, errors.New("Failed to read game log file(Oversea)")
	}

	regexpGameDataDir, err := regexp.Compile(exprGameDataDir)
	if err != nil {
		return nil, errors.New("Failed to compile exprGameDataDir(Oversea)")
	}
	resultGameDataDir := regexpGameDataDir.FindSubmatch(logDataBytes)
	if resultGameDataDir == nil {
		return nil, errors.New("Failed to find game data directory(Oversea)")
	}
	gameDataDir := filepath.Join(string(resultGameDataDir[1]), "LocalCache/Data")

	regexpGachaRecordUrl, err := regexp.Compile(exprGachaRecordUrlOversea)
	if err != nil {
		return nil, errors.New("Failed to compile exprGachaRecordUrl(Oversea)")
	}
	resultGachaRecordUrlList := regexpGachaRecordUrl.FindAllSubmatch(logDataBytes, -1)
	if resultGachaRecordUrlList == nil {
		return nil, errors.New("Failed to find gacha record url(Oversea)")
	}
	latestResultGachaRecordUrl := resultGachaRecordUrlList[len(resultGachaRecordUrlList)-1]
	gachaRecordUrl := string(latestResultGachaRecordUrl[1])

	regexpLoginInfo, err := regexp.Compile(exprLoginInfoOversea)
	if err != nil {
		return nil, errors.New("Failed to compile exprLoginInfo(Oversea)")
	}
	resultLoginInfoList := regexpLoginInfo.FindAllSubmatch(logDataBytes, -1)
	if resultLoginInfoList == nil {
		return nil, errors.New("Failed to find game login information(Oversea)")
	}
	latestResultLoginInfo := resultLoginInfoList[len(resultLoginInfoList)-1]
	gameAccessToken := string(latestResultLoginInfo[2])
	uid, err := strconv.ParseUint(string(latestResultLoginInfo[1]), 10, 64)
	if err != nil {
		return nil, errors.Errorf("Failed to parse uid:'%s'(Oversea)", latestResultLoginInfo[1])
	}

	var server gameServer
	switch {
	case strings.Contains(gachaRecordUrl, "gf2-gacha-record-us"):
		server = GameServerUS
	case strings.Contains(gachaRecordUrl, "gf2-gacha-record-intl"):
		server = GameServerGlobal
	case strings.Contains(gachaRecordUrl, "gf2-gacha-record-jp"):
		server = GameServerJP
	case strings.Contains(gachaRecordUrl, "gf2-gacha-record-kr"):
		server = GameServerKR
	case strings.Contains(gachaRecordUrl, "gf2-gacha-record-asia"):
		server = GameServerAsia
	default:
		return nil, errors.Errorf("Failed to determine server(Oversea),gacha url:%s", gachaRecordUrl)
	}

	return &GameUserInfo{
		IdentitySource:  "gameLog",
		Uid:             uid,
		GameServer:      server,
		GameDataDir:     gameDataDir,
		GameAccessToken: gameAccessToken,
		GachaRecordUrl:  gachaRecordUrl,
	}, nil
}
