package download

import (
	db "Octave/golibs/database"
	"Octave/golibs/settings"
	"fmt"

	ytsearch "github.com/AnjanaMadu/YTSearch"
	youtube "github.com/knadh/go-get-youtube/youtube"
)

func Download(song db.Song) error {

	results, err := ytsearch.Search(fmt.Sprintf("%s %s", song.Artist, song.Title))
	if err != nil {
		return err
	}
	first := results[0]
	vidid := first.VideoId
	loc := settings.Settings().SongDir + song.ID + ".mp3"
	err = fromID(vidid, loc)
	if err != nil {

		return err
	}

	return nil
}

func fromID(id string, location string) error {

	// get the video object (with metdata)
	video, err := youtube.Get(id)

	if err != nil {
		return err
	}
	// download the video and write to file
	option := &youtube.Option{
		Rename: false, // rename file using video title
		Resume: true,  // resume cancelled download
		Mp3:    true,  // extract audio to MP3
	}

	return video.Download(0, location, option)
}
