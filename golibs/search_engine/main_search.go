package searchengine

import (
	db "Octave/golibs/database"
	"Octave/golibs/settings"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/meilisearch/meilisearch-go"
	ps "github.com/mitchellh/go-ps"
)

var (
	client = meilisearch.NewClient(meilisearch.ClientConfig{
		Host: "http://127.0.0.1:7700",
	})
	ind = client.Index("songs")
)
var results []db.Song

func StartEngine() error {
	isrunning := false
	processList, err := ps.Processes()
	if err != nil {
		return err
	}
	for _, process := range processList {
		if process.Executable() == "OctaveSearchEngine.exe" {
			isrunning = true

			break
		}
	}
	if !isrunning {
		sett := settings.Settings()
		go exec.Command("./OctaveSearchEngine.exe", "--db-path", "./meilli", //nolint:gosec
			"--max-indexing-memory", fmt.Sprint(sett.MeilliRAM),
			"--max-indexing-threads", fmt.Sprint(sett.MeilliThreads),
			"--http-addr", fmt.Sprintf("127.0.0.1:%v", sett.MeilliPort)).Start() //nolint

	}

	return nil

}

func FirstIndex() error {
	js, err := json.Marshal(db.OpenDatabase().GetAllSongs())
	if err != nil {

		return err
	}
	client := meilisearch.NewClient(meilisearch.ClientConfig{
		Host: "http://127.0.0.1:7700",
	})
	_, err = client.Index("songs").AddDocuments(js)
	if err != nil {

		return err
	}

	return nil
}

func Search(query string) ([]db.Song, error) {
	results = []db.Song{}
	var searchRequest meilisearch.SearchRequest
	start := time.Now()
	res, err := ind.Search(query, &searchRequest)
	fmt.Printf("Search took %s\n", time.Since(start))
	if err != nil {

		return []db.Song{}, err
	}
	for _, i := range res.Hits {
		result := db.Song{
			Title:  i.(map[string]any)["Title"].(string),
			Artist: i.(map[string]any)["Artist"].(string),
			Album:  i.(map[string]any)["Artist"].(string),
			Length: int(i.(map[string]any)["Length"].(float64)),
			ID:     i.(map[string]any)["id"].(string),
			Image:  "",
		}

		if i.(map[string]any)["Image"] != nil {
			result.Image = i.(map[string]any)["Image"].(string)
		} else {
			result.Image = ""
		}
		if err != nil {

			return []db.Song{}, err
		}

		results = append(results, result)

	}

	return results, nil

}
