package settings

import (
	"encoding/json"
	"log"
	"os"
)

var settingsmap = make(map[string]any)

var set1 []byte

type SettingStruct struct {
	Version       string `json:"version"`       // Version of the application
	Name          string `json:"name"`          // Name of the application (Probably won't change)
	User          string `json:"user"`          // Username of the user. Not used right now.
	SongDir       string `json:"songDir"`       // Directory where songs are stored
	DBPath        string `json:"dbPath"`        // Path to the database
	WSport        int    `json:"wsport"`        // Port for the websockets
	LogPath       string `json:"logPath"`       // Path to the logs folder
	StateFile     string `json:"stateFile"`     // Json file to store the user's state
	MusicdbToken  string `json:"musicdbToken"`  // Token for the musicdb. Not used
	MeilliThreads int    `json:"meilliThreads"` // Number of threads for the Meilli server
	MeilliRAM     int    `json:"meilliRam"`     // Amount of RAM for the Meilli server
	MeilliPort    int    `json:"meilliPort"`    // Port to run the search engine on
	AssetDir      string `json:"assetDir"`      // Where assets are served from when not compiled with the app
}

func SetEmbedded(settings []byte) {
	set1 = settings
}
func Settings() SettingStruct {

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
	sett := SettingStruct{}
	sett.Version = settingsmap["version"].(string)
	sett.Name = settingsmap["name"].(string)
	sett.User = settingsmap["user"].(string)
	sett.SongDir = settingsmap["songDir"].(string)
	sett.DBPath = settingsmap["dbPath"].(string)
	settingsmap["wsport"] = int(settingsmap["wsport"].(float64))
	sett.WSport = int(settingsmap["wsport"].(int))
	sett.LogPath = settingsmap["logPath"].(string)
	sett.StateFile = settingsmap["stateFile"].(string)
	sett.MusicdbToken = settingsmap["musicdbToken"].(string)
	sett.MeilliThreads = int(settingsmap["meilliThreads"].(float64))
	sett.MeilliRAM = int(settingsmap["meilliRam"].(float64))
	sett.MeilliPort = int(settingsmap["meilliPort"].(float64))
	sett.AssetDir = settingsmap["assetDir"].(string)

	return sett
}
