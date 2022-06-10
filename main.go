package main

import (
	db "Octave/golibs/database"
	logging "Octave/golibs/log"
	"Octave/golibs/search_engine"
	"Octave/golibs/settings"
	global_state "Octave/golibs/state"
	"embed"
	"encoding/json"
	"flag"
	"os"
	"runtime/debug"

	golog "github.com/apsdehal/go-logger"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/panicwrap"
	"github.com/wailsapp/wails/v2"
	wails_logger "github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	_ "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed frontend/src
var Assets embed.FS

var Settings settings.SettingsStruct
var Log *golog.Logger

func fatalpanic(stacktrace string) {
	// This is pretty much a final step for if the app panics. We don't try to recover cos that won't work
	// Instead we just record as much data as we can
	Log.Warning("Shit just hit the fan in a way we can't recover from")
	Log.Warning("Stack trace below")
	Log.Criticalf("%s", stacktrace)
	Log.Warning("Freeing up memory before quitting")
	debug.FreeOSMemory()
	logging.Close(Log)
}

//go:embed "settings.json"
var set1 []byte

func main() {

	settings.SetEmbedded(set1)
	logging.CreateLogger()
	Log = logging.Log

	// Read command line flags
	catchpanic := flag.Bool("debug", false, "A bool")
	flag.Parse()
	// Only run the panic wrapper if we're not in debug mode, cos it messes with vscode's debugging
	// so just use -debug as a flag to disable it, which I added to the debug profile
	if *catchpanic {
		println("Running in production mode")
		// When the app closes from either user input or just stuff breaking, run this function
		defer func() {
			logging.Close(Log)
		}()
		exitStatus, err := panicwrap.BasicWrap(fatalpanic)
		if err != nil {
			// Something went wrong setting up the panic wrapper. Unlikely,
			// but possible.
			Log.FatalF("So. The system designed to catch errors, has errored.\n%s", err)
		}

		// If exitStatus >= 0, then we're the parent process and the panicwrap
		// re-executed ourselves and completed. Just exit with the proper status.
		if exitStatus >= 0 {
			os.Exit(exitStatus)
		}
	}
	go search_engine.StartEngine()
	go search_engine.FirstIndex()
	Settings = settings.Settings()
	// Create an instance of the app structure
	Log.Info("Creating App")
	app := NewApp()
	var state global_state.Statemap
	// Open Settings.stateFile and read the json into a statemap
	Log.Info("Loading StateFile")
	statefile, err := os.ReadFile(Settings.StateFile)
	if err != nil {
		Log.Error(err.Error())
	}
	Log.Info("Unmarshalling StateFile into Statemap")
	err = json.Unmarshal(statefile, &state)
	if err != nil {
		Log.Error(err.Error())
	}
	Log.Info("Updating global_state.State")
	global_state.SaveState(state)
	Log.Info("Setting up websocket")
	go handleRequests()
	var songs []db.Song
	allsongs := db.OpenDatabase().GetAllSongs()
	songs = append(songs, allsongs...)
	global_state.State.ActivePlaylist = db.Playlist{
		ID:    "0",
		Name:  "Default",
		Songs: songs,
		Art:   "witcher.png",
	}
	// Create application with options
	Log.Info("Creating Wails App")
	err = wails.Run(&options.App{
		Title: Settings.Name,
		//Width:             1024,
		//Height:            768,
		DisableResize:     true,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		RGBA:              &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            Assets,
		Menu:              nil,
		Logger:            &wails_logger.DefaultLogger{},
		LogLevel:          wails_logger.TRACE,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnBeforeClose:     app.beforeClose,
		OnShutdown:        app.shutdown,
		WindowStartState:  options.Fullscreen,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
		},
	})

	if err != nil {
		Log.ErrorF(err.Error())
	}
}
