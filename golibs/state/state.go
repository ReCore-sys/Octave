package global_state

import (
	db "Octave/golibs/database"
	logging "Octave/golibs/log"
)

type Statemap struct {
	CurrentSong    db.Song     `json:"currentSong"`    // the ID of the current song
	Elapsed        float64     `json:"elapsed"`        // How many seconds into the current song we are
	Paused         bool        `json:"paused"`         // If the player is paused
	Volume         float64     `json:"volume"`         // The volume of the current song
	Muted          bool        `json:"muted"`          // If the player is muted
	Repeat         bool        `json:"repeat"`         // If repeat mode is active
	Shuffle        bool        `json:"shuffle"`        // If the player is on shuffle mode
	Next           db.Song     `json:"next"`           // The ID of the next song
	Prev           db.Song     `json:"prev"`           // The ID of the previous song
	ActivePlaylist db.Playlist `json:"activePlaylist"` // The ID of the active playlist
	//Queue          []db.Song   `json:"queue"`          // The IDs of the songs in the queue
	CurrentIndex int `json:"currentIndex"` // The index of the current song in the queue
}

var (
	State Statemap
)

func SaveState(s Statemap) {
	Log := logging.Log
	Log.Info("Saving state")
	State = s
}
