package settings

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"
)

var settingsmap = make(map[string]any)

var set1 []byte

type SettingsStruct struct {
	Version       string `json:"version"`
	Name          string `json:"name"`
	User          string `json:"user"`
	SongDir       string `json:"songDir"`
	DBPath        string `json:"dbPath"`
	WSport        int    `json:"wsport"`
	LogPath       string `json:"logPath"`
	StateFile     string `json:"stateFile"`
	Musicdb_token string `json:"musicdb_token"`
	MeilliThreads int    `json:"MeilliThreads"`
	MeilliRam     int    `json:"MeilliRam"`
	MeilliPort    int    `json:"MeilliPort"`
	AssetDir      string `json:"AssetDir"`
}

func SetEmbedded(settings []byte) {
	set1 = settings
}
func Settings() SettingsStruct {

	set2, err := os.ReadFile("user-settings.json")
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(set1, &settingsmap)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(set2, &settingsmap)
	if err != nil {
		log.Println(err)
	}
	sett := SettingsStruct{}
	sett.Version = settingsmap["version"].(string)
	sett.Name = settingsmap["Name"].(string)
	sett.User = settingsmap["User"].(string)
	sett.SongDir = settingsmap["SongDir"].(string)
	sett.DBPath = settingsmap["dbPath"].(string)
	settingsmap["wsport"] = int(settingsmap["wsport"].(float64))
	sett.WSport = int(settingsmap["wsport"].(int))
	sett.LogPath = settingsmap["logPath"].(string)
	sett.StateFile = settingsmap["stateFile"].(string)
	sett.Musicdb_token = settingsmap["musicdb_token"].(string)
	sett.MeilliThreads = int(settingsmap["MeilliThreads"].(float64))
	sett.MeilliRam = int(settingsmap["MeilliRam"].(float64))
	sett.MeilliPort = int(settingsmap["MeilliPort"].(float64))
	sett.AssetDir = settingsmap["AssetDir"].(string)
	return sett
}
