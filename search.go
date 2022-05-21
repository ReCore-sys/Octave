package main

import (
	db "Octave/golibs/database"
	se "Octave/golibs/search_engine"
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

var cache []db.Song

func (a *App) Search(input string) []db.Song {
	var results []db.Song
	if input == "test" {
		dummysong := []db.Song{
			{
				ID:     "0",
				Title:  "Another one bites the dust",
				Artist: "Queen",
				Album:  "The Game",
				Length: 224,
				Image:  "https://upload.wikimedia.org/wikipedia/en/4/4d/Another_one_bites_the_dust.jpg",
			},
		}
		cache = dummysong
		return dummysong
	}

	Log.Infof("Searching for: %v", input)
	urlparts := url.Values{}
	urlparts.Set("q", input)
	qurl := "https://api.deezer.com/search?" + urlparts.Encode()
	println(qurl)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", qurl, nil)
	if err != nil {
		Log.Errorf("Error creating request: %s", err)
		return results
	}
	req.Header.Set("User-Agent", "Octave/01")
	resp, err := client.Do(req)
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
		song := db.Song{}

		if InMap("id", result) {
			song.ID = strconv.Itoa(int(result["id"].(float64)))
		} else {
			continue
		}

		if InMap("title", result) {
			song.Title = result["title"].(string)
		} else {
			continue
		}

		if InMap("artist", result) {
			song.Artist = result["artist"].(map[string]any)["name"].(string)
		} else {
			continue
		}

		if InMap("album", result) {
			song.Album = result["album"].(map[string]any)["title"].(string)
		} else {
			continue
		}

		if InMap("duration", result) {
			song.Length = int(result["duration"].(float64))
		} else {
			continue
		}
		defer func() {
			if r := recover(); r != nil {
				song.Image = ""
			}
		}()
		if InMap("album", result) {
			sizes := []string{"xl", "big", "medium", "small"}
			for _, size := range sizes {
				if InMap("cover_"+size, result["album"].(map[string]any)) {
					if reflect.TypeOf(result["album"].(map[string]any)["cover_"+size]) != nil {

						song.Image = result["album"].(map[string]any)["cover_"+size].(string)
						break
					} else {
						Log.WarningF("Song %s had a nil image", result["title"].(string))
					}
				}
			}
		}

		results = append(results, song)
	}
	Log.Info("Search complete")
	cache = results
	return results

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
