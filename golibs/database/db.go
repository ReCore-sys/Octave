package db

import (
	logging "Octave/golibs/log"
	"Octave/golibs/settings"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type DB struct {
	db *sql.DB
}

var Log = logging.Log

func GenerateID(path string) (string, error) {
	settings := settings.Settings()
	fullpath := settings.SongDir + path
	// Hash the file
	file, err := ioutil.ReadFile(fullpath)
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	hash.Write(file)
	return string(hash.Sum(nil)), nil

}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func CreateString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)

}

func OpenDatabase() DB {
	if Log == nil {
		logging.CreateLogger()
		Log = logging.Log
	}
	db, err := sql.Open("sqlite3", "./music.db")
	if err != nil {
		log.Println(err)
	}
	return DB{db: db}
}

type Playlist struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Songs []Song `json:"songs"`
	Art   string `json:"art"`
	Color string `json:"color"`
}

type Song struct {
	ID     string `json:"id"`     // Unique ID of the song
	Title  string `json:"Title"`  // Display name of the song
	Artist string `json:"Artist"` // Artist of the song
	Album  string `json:"Album"`  // What album the song is from
	Length int    `json:"Length"` // Length of the song in seconds
	Image  string `json:"Image"`  // Path to the image for the song
}

func (db DB) LookupSong(id string) Song {
	row := db.db.QueryRow("SELECT * FROM songs WHERE ID = ?", id)
	var title, artist, album string
	var length int
	var image string
	err := row.Scan(&title, &album, &artist, &length, &image, &id)
	if err != nil {

		Log.Infof("Looking for song %s", id)
		formatted := fmt.Sprint(err)
		Log.Error(formatted)
	}
	return Song{
		Title:  title,
		Artist: artist,
		Album:  album,
		Length: length,
		Image:  image,
		ID:     id,
	}

}

func (db DB) GetAllSongs() []Song {
	rows, err := db.db.Query("SELECT * FROM songs")
	if err != nil {
		log.Println(err)
	}
	var songs []Song
	for rows.Next() {
		var title, artist, album, image, id string
		var length int
		err := rows.Scan(&title, &album, &artist, &length, &image, &id)
		if err != nil {
			log.Println(err)
		}
		songs = append(songs, Song{
			ID:     id,
			Title:  title,
			Artist: artist,
			Album:  album,
			Length: length,
			Image:  image,
		})
	}
	return songs
}

func (db DB) AddSong(song Song) error {
	_, err := db.db.Exec("INSERT INTO songs VALUES (?, ?, ?, ?, ?,  ?)", song.Title, song.Album, song.Artist, song.Length, song.Image, song.ID)
	return err
}

func (db DB) DoesSongExist(id string) bool {
	var exists bool
	row := db.db.QueryRow("SELECT EXISTS(SELECT 1 FROM songs WHERE ID = ?)", id)
	err := row.Scan(&exists)
	if err != nil {
		log.Println(err)
	}
	return exists
}

/**------------------------------------------------------------------------
 **                            Playlist Functions
 *------------------------------------------------------------------------**/

// Creating a playlist.
func (db DB) CreatePlaylist(name string) (bool, Playlist) {
	var id string
	tx, err := db.db.Begin()
	if err != nil {
		Log.Error(err.Error())
		return false, Playlist{}
	}
	for {
		id = CreateString(32)
		var exists bool
		row := db.db.QueryRow("SELECT EXISTS(SELECT 1 FROM playlists WHERE ID = ?)", id)
		err := row.Scan(&exists)
		if err != nil {
			Log.Error(err.Error())
			return false, Playlist{}
		}
		if !exists {
			break
		}

	}
	_, err = tx.Exec("INSERT INTO playlists VALUES (?, ?, \"#0\", 0)", id, name)
	if err != nil {
		tx.Rollback()
		Log.Error(err.Error())
		return false, Playlist{}
	}
	_, err = tx.Exec("CREATE TABLE playlist_" + id + " (song_id TEXT, order_ int, art TEXT)")
	if err != nil {
		tx.Rollback()
		Log.Error(err.Error())
		return false, Playlist{}
	}
	tx.Commit()
	return db.GetPlaylist(id)
}

// A function that takes a DB and a string and returns a bool and a Playlist.
func (db DB) GetPlaylist(id string) (bool, Playlist) {
	/*
		Playlists record table format:
			Playlist ID, Playlist Name, Playlist Length

		Individual Playlist Table Format:
			Playlist ID, Song ID, Order
	*/
	rows, err := db.db.Query("SELECT name,color FROM playlists WHERE ID = ?", id)
	if err != nil {
		Log.Error(err.Error())
		return false, Playlist{}
	}
	var name string
	var color string
	for rows.Next() {
		err := rows.Scan(&name, &color)
		if err != nil {
			Log.Error(err.Error())
			return false, Playlist{}
		}
	}
	rows, err = db.db.Query("SELECT song_id, order_, art FROM playlist_" + id)
	if err != nil {
		Log.Error(err.Error())
		return false, Playlist{}
	}
	var songs []Song
	var art string
	for rows.Next() {
		var song string
		var order int
		err := rows.Scan(&song, &order, &art)
		if err != nil {
			Log.Error(err.Error())
			return false, Playlist{}
		}
		truesong := db.LookupSong(song)
		songs = append(songs, truesong)
	}
	return true, Playlist{
		ID:    id,
		Name:  name,
		Songs: songs,
		Art:   art,
		Color: color,
	}

}

func (db DB) GetAllPlaylists() ([]Playlist, error) {
	rows, err := db.db.Query("SELECT id FROM playlists")
	if err != nil {
		Log.Error(err.Error())
		return nil, err
	}
	var playlists []Playlist
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			Log.Error(err.Error())
			return nil, err
		}
		exists, pl := db.GetPlaylist(id)
		if !exists {
			return nil, errors.New("Playlist not found")
		}
		playlists = append(playlists, pl)
	}
	return playlists, nil
}

// Deleting a playlist from the database.
func (db DB) DeletePlaylist(id string) {
	db.db.Exec("DROP TABLE playlist_" + id)
	db.db.Exec("DELETE FROM playlists WHERE ID = ?", id)
}
