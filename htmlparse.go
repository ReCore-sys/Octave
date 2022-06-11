package main

import (
	db "Octave/golibs/database"
	logging "Octave/golibs/log"
	"Octave/golibs/settings"
	global_state "Octave/golibs/state"
	"encoding/json"

	"github.com/flosch/pongo2"
)

func (a *App) Parse(path string) string {
	Settings := settings.Settings()
	if Log == nil {
		logging.CreateLogger()
		Log = logging.Log
	}
	Log.Info("Parsing template: " + path)
	file, err := Assets.ReadFile("frontend/src/templates/" + path)
	if err != nil {
		Log.Error(err.Error())
	}
	t, err := pongo2.FromBytes(file)
	if err != nil {
		Log.Error(err.Error())
	}

	// Some paths require certain variables
	specialmap := make(map[string]any)
	if path == "playlists_sidebar.html" {
		playlists, err := db.OpenDatabase().GetAllPlaylists()
		if err != nil {
			Log.Error(err.Error())
		}
		specialmap["playlists"] = playlists
	}

	maps := MapMixer(StructToMap(Settings), StructToMap(global_state.State), specialmap)
	var emptycontext pongo2.Context
	v, err := json.Marshal(maps)
	if err != nil {
		Log.Error(err.Error())
	}
	err = json.Unmarshal(v, &emptycontext)
	if err != nil {
		Log.Error(err.Error())
	}
	out, err := t.Execute(emptycontext)
	if err != nil {
		Log.Error(err.Error())
	}

	return out
}

func MapMixer(inmaps ...map[string]any) map[string]any {
	outmap := make(map[string]any)
	for _, inmap := range inmaps {
		for k, v := range inmap {
			outmap[k] = v
		}
	}

	return outmap
}

func StructToMap(str any) map[string]any {
	outmap := make(map[string]any)
	jsonstr, err := json.Marshal(str)
	if err != nil {
		Log.Error(err.Error())
	}
	err = json.Unmarshal(jsonstr, &outmap)
	if err != nil {
		Log.Error(err.Error())
	}

	return outmap
}
