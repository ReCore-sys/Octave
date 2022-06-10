package main

import (
	db "Octave/golibs/database"
	"Octave/golibs/settings"
	global_state "Octave/golibs/state"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	state := global_state.State
	// Convert to json then write to file
	loc := Settings.StateFile
	json, err := json.Marshal(state)
	if err != nil {
		Log.Errorf("Error marshalling state: %s", err)
	}
	err = ioutil.WriteFile(loc, json, 0644)
	if err != nil {
		Log.Errorf("Error writing state: %s", err)
	}

	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

/* Overwriting the index.html file with the contents of the trueindex.html file. */
//func (a *App) Overwrite() string {
//	text := a.Parse("trueindex.html")
//	return text
//}

/* It's returning the global state. */
func (a *App) GetState() global_state.Statemap {
	return global_state.GetState()
}

/* It's updating the current song in the global state. */
func (a *App) UpdateSong(song db.Song) {
	state := global_state.State
	state.CurrentSong = song
	global_state.SaveState(state)
}

func (a *App) FindSong(id string) db.Song {
	database := db.OpenDatabase()
	return database.LookupSong(id)
}

func (a *App) GetAllSongs() []db.Song {
	database := db.OpenDatabase()
	return database.GetAllSongs()
}

func (a *App) Sleep() string {
	Log.Debug("Sleeping")
	time.Sleep(time.Second * 5)
	Log.Debug("Waking up")
	return "Done"
}

/* It's a function that takes a string and a logtype and logs it to the console. */
func (a *App) JSLog(logtype string, msg string) {
	switch logtype {
	case "error":
		Log.Errorf("JS error: %v", msg)
	case "info":
		Log.Infof("JS Info: %v", msg)
	case "debug":
		Log.Debugf("JS Debug: %v", msg)
	case "fatal":
		Log.Fatalf("JS Fatal: %v", msg)
	case "warn":
		Log.Warningf("JS Warn: %v", msg)
	}
}

func (a *App) Save(id string) bool {
	Log.Infof("Saving song %s", id)
	database := db.OpenDatabase()
	song := a.InCache(id)
	go SaveImage(song)
	song.Image = fmt.Sprintf("%s.jpg", id)
	if (song == db.Song{}) {
		Log.ErrorF("Song not found in cache: %v", id)
		return false
	}
	err := database.AddSong(song)
	return err == nil
}

func (a *App) SongDownloaded(id string) bool {
	database := db.OpenDatabase()
	ex := database.DoesSongExist(id)
	if ex {
		Log.Infof("Song %s already exists", id)
	}
	return ex
}

func (a *App) Settings() settings.SettingsStruct {
	if (Settings == settings.SettingsStruct{}) {
		Settings = settings.Settings()
	}
	return Settings
}

func (a *App) CreatePlaylist(name string) bool {
	database := db.OpenDatabase()
	err, _ := database.CreatePlaylist(name)
	return err
}

func SaveImage(song db.Song) {
	Log.Infof("Saving image for %s", song.ID)
	resp, err := http.Get(song.Image)
	if err != nil {
		Log.Errorf("Error getting image: %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Errorf("Error reading image: %s", err)
		return
	}
	loc := fmt.Sprintf("%s/song_img/%s.jpg", Settings.AssetDir, song.ID)
	f, err := os.OpenFile(loc, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Log.Errorf("Error opening file: %s", err)
		return
	}
	defer f.Close()
	f.Write(body)

}
