package main

import (
	db "Octave/golibs/database"
	"Octave/golibs/download"
	logging "Octave/golibs/log"
	"Octave/golibs/search_engine"
	"Octave/golibs/settings"
	global_state "Octave/golibs/state"
	_ "embed"
	"encoding/json"
	"os"
	"testing"
	"time"
)

var nums = []int{
	84, 76, 15, 92, 74, 23, 76, 55, 33, 17, 79, 75, 27, 5, 78, 58, 84, 69, 45, 43, 38, 24, 65, 65, 51, 97, 18, 51, 83, 11, 84, 60, 69, 15, 1, 96, 90, 59, 81, 65, 27, 1, 40, 46, 7, 99, 54, 72, 92, 77, 92, 55, 86, 88, 54, 42, 5, 12, 91, 16, 79, 52, 68, 31, 77, 80, 34, 17, 81, 58, 92, 49, 22, 46, 53, 36, 39, 28, 88, 10, 82, 42, 70, 13, 54, 88, 69, 79, 96, 64, 28, 24, 20, 12, 33, 12, 91, 26, 15, 8, 46, 40, 26, 12, 21, 23, 87, 76, 94, 72, 9, 89, 26, 47, 1, 79, 90, 68, 60, 48, 64, 90, 45, 6, 38, 46, 65, 44, 85, 20, 70, 78, 95, 55, 34, 68, 72, 52, 53, 68, 75, 16, 69, 37, 89, 53, 8, 41, 4, 87, 95, 74, 62, 31, 31, 18, 6, 83, 79, 28, 84, 66, 30, 59, 85, 12, 93, 37, 42, 52, 39, 85, 46, 76, 90, 29, 90, 2, 22, 66, 50, 52, 50, 22, 66, 92, 48, 72, 2, 27, 37, 90, 83, 77, 52, 90, 48, 46, 95, 25, 54, 30, 60, 23, 97, 22, 79, 76, 91, 10, 41, 85, 82, 7, 2, 59, 62, 39, 44, 65, 49, 55, 59, 25, 67, 70, 29, 17, 29, 85, 18, 71, 75, 74, 59, 59, 96, 17, 70, 8, 28, 25, 85, 74, 87, 48, 83, 52, 25, 37, 98, 1, 20, 35, 44, 3, 39, 15, 79, 7, 51, 72, 78, 15, 52, 11, 16, 4, 16, 38, 22, 53, 9, 80, 64, 77, 29, 84, 16, 35, 21, 87, 29, 45, 41, 2, 20, 52, 92, 85, 18, 69, 75, 19, 7, 89, 36, 35, 55, 59, 16, 81, 38, 57, 25, 80, 96, 16, 82, 31, 97, 94, 20, 93, 78, 76, 41, 30, 83, 16, 89, 24, 37, 82, 36, 40, 4, 46, 46, 33, 46, 6, 28, 31, 29, 54, 96, 46, 26, 37, 56, 10, 27, 84, 88, 5, 31, 66, 45, 29, 87, 85, 52, 36, 18, 95, 2, 41, 56, 60, 98, 99, 61, 79, 77, 58, 41, 40, 40, 3, 90, 93, 49, 83, 80, 55, 96, 50, 92, 25, 47, 69, 16, 50, 23, 23, 63, 27, 42, 50, 36, 72, 94, 95, 94, 60, 20, 21, 59, 57,
}

var target = 55

func BenchmarkFindCon(b *testing.B) {
	var answer int
	found := false
	findfunc := func(inp int, index int) {
		if found == false {
			if nums[index] == inp {
				answer = index
				found = true
			}
		}

	}
	for ind := range nums {
		if found == false {
			go findfunc(target, ind)
		} else {
			break
		}
	}
	time.Sleep(50 * time.Millisecond)
	if !found {
		b.Error("not found")
	} else {
		b.Log(answer)
	}
}
func BenchmarkFind(b *testing.B) {
	for i, val := range nums {
		if val == target {
			b.Log(i)
			break
		}
	}
}

func BenchmarkSearch(b *testing.B) {
	q := "Stains of time"
	a := App{}
	res := a.Search(q)
	b.Logf("%v: %v", res[0].Title, res[0].ID)
}

func BenchmarkSearchDisk(b *testing.B) {
	q := "Stains of time"
	res, err := search_engine.Search(q)
	if err != nil {
		b.Error(err)
	}
	b.Logf("%v: %v", res[0].Title, res[0].ID)
}

func BenchmarkIndex(b *testing.B) {
	search_engine.FirstIndex()

}

func BenchmarkTemplate(b *testing.B) {
	a := App{}
	a.Parse("home.html")

}

//go:embed "settings.json"
var default_settings []byte

func TestDownload(t *testing.T) {
	settings.SetEmbedded(default_settings)
	song := db.OpenDatabase().LookupSong("1104838")
	err := download.Download(song)
	if err != nil {
		t.Error(err)

	}
}

func TestHTML(t *testing.T) {

	settings.SetEmbedded(default_settings)
	if Log == nil {
		logging.CreateLogger()
		Log = logging.Log
	}
	Settings = settings.Settings()
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
	a := App{}
	b := a.Parse("trueindex.html")
	t.Log(b)
}

func TestState(t *testing.T) {
	settings.SetEmbedded(default_settings)
	if Log == nil {
		logging.CreateLogger()
		Log = logging.Log
	}
	Settings = settings.Settings()
	stateFile, err := os.Open(Settings.StateFile)
	if err != nil {
		t.Error(err)

	}
	var emptystate global_state.Statemap
	err = json.NewDecoder(stateFile).Decode(&emptystate)
	if err != nil {
		t.Errorf("Error decoding 1: %v", err)
	}
	gs := global_state.State
	jsonform, err := json.Marshal(gs)
	if err != nil {
		t.Errorf("Error marshalling 2: %v", err)
	}
	var marshaltarget global_state.Statemap
	err = json.Unmarshal(jsonform, &marshaltarget)
	if err != nil {
		t.Errorf("Error unmarshalling 3: %v", err)
	}

}
