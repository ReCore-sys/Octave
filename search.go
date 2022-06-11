package main

import (
	db "Octave/golibs/database"
	se "Octave/golibs/search_engine"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
)

var cache []db.Song

var results []db.Song

func (a *App) Search(input string) []db.Song {
	results = []db.Song{}

	Log.Infof("Searching for: %v", input)
	urlparts := url.Values{}
	urlparts.Set("q", input)
	resp, err := http.Get("https://api.deezer.com/search?" + urlparts.Encode())
	if err != nil {
		Log.Error(err.Error())
	}
	defer resp.Body.Close()
	resmap := make(map[string]any)
	err = json.NewDecoder(resp.Body).Decode(&resmap)
	if err != nil {
		Log.Error(err.Error())
	}
	resultlist := resmap["data"].([]any)
	for _, r := range resultlist {
		result := r.(map[string]any)
		song, ok := ResultsToSong(result)
		if !ok {
			continue
		}

		results = append(results, song)
	}
	Log.Info("Search complete")
	cache = results

	return results

}

// Given the json encoding of a result from the API, this function converts it to a song.
func ResultsToSong(result map[string]any) (db.Song, bool) {

	song := db.Song{}
	if InMap("id", result) {
		song.ID = strconv.Itoa(int(result["id"].(float64)))
	} else {
		return song, false
	}

	if InMap("title", result) {
		song.Title = Sanitize(result["title"].(string))
	} else {
		return song, false
	}

	if InMap("artist", result) {
		song.Artist = Sanitize(result["artist"].(map[string]any)["name"].(string))
	} else {
		return song, false
	}

	if InMap("album", result) {
		song.Album = Sanitize(result["album"].(map[string]any)["title"].(string))
	} else {
		return song, false
	}

	if InMap("duration", result) {
		song.Length = int(result["duration"].(float64))
	} else {
		return song, false
	}
	if InMap("album", result) {
		sizes := []string{"xl", "big", "medium", "small"}
		for _, size := range sizes {
			if InMap("cover_"+size, result["album"].(map[string]any)) {
				if reflect.TypeOf(result["album"].(map[string]any)["cover_"+size]) != nil {

					song.Image = Sanitize(result["album"].(map[string]any)["cover_"+size].(string))

					break
				} else {
					Log.WarningF("Song %s had a nil image", result["title"].(string))
				}
			}
		}
	}

	return song, true

}

// When you run the search command, it is saved in the cache. This function allows you to retrieve a result by ID.
func (a *App) InCache(id string) db.Song {
	for _, result := range cache {
		if result.ID == id {

			return result
		}
	}

	return db.Song{}
}

/**--------------------------------------------
 *               PAIN.JPG
 *---------------------------------------------**/

// "Do a search engine, it'll be cool" I thought. It is not cool and I hate my life

// Searching the database for a song that matches the query.
func (a *App) SearchSaved(query string) []db.Song {
	res, err := se.Search(query)
	if err != nil {
		Log.Error(err.Error())
	}
	cache = res

	return res

}

func InMap(item string, mep map[string]any) bool {
	_, ok := mep[item]

	return ok

}

func RegexMatch(query string, target string) bool {
	r, err := regexp.Compile(query)
	if err != nil {
		Log.Error(err.Error())

		return false
	}

	return r.MatchString(target)
}

func Sanitize(s string) string {
	regmatch := "<script>(.*?)</script>"
	if RegexMatch(regmatch, s) {
		innerscript, err := regexp.Compile(regmatch)
		if err != nil {
			Log.Error(err.Error())
		}
		s = string(innerscript.FindSubmatch([]byte(s))[0])
	}

	return s
}
