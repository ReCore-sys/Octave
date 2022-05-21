package main

import (
	logging "Octave/golibs/log"
	global_state "Octave/golibs/state"
	"bytes"
	"encoding/json"
	"text/template"
)

func (a *App) Parse(path string) string {
	if Log == nil {
		logging.CreateLogger()
		Log = logging.Log
	}
	Log.Info("Parsing template: " + path)
	file, err := Assets.ReadFile("frontend/src/templates/" + path)
	if err != nil {
		Log.Error(err.Error())
	}
	Log.Info("Parsing template: " + path + " complete")
	filetext := string(file)
	t := template.Must(template.New("main").Parse(filetext))
	// Create an empty io reader and writer

	var final bytes.Buffer
	Log.Info("Mixing maps")
	maps := MapMixer(StructToMap(Settings), StructToMap(global_state.State))
	err = t.ExecuteTemplate(&final, "main", maps)
	if err != nil {
		Log.Error(err.Error())
	}
	cont := final.String()

	// find any instance of {{ import <file> }} and replace with the contents of the file
	// this is a hacky way to do this, but it works for now

	return cont
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
