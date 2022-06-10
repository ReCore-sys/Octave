package logging

import (
	"Octave/golibs/settings"
	"io"
	"log"
	"os"
	"strings"

	"time"

	"github.com/apsdehal/go-logger"
)

var Log *logger.Logger
var created = false

/* Creating a logger. */
func CreateLogger() {
	if created {
		Log.Warning("Logger already created")
		return
	}
	var err error
	fmat := "%{id} %{time} %{file}:%{line} | [%{level}] >  %{message}"
	sett := settings.Settings()
	// Use layout string for time format.
	const layout = "01-02-2006"
	// Place now in the string.
	t := time.Now()
	loc := sett.LogPath + "/" + t.Format(layout) + ".log"
	latestname := sett.LogPath + "/latest.log"
	datefile, err := os.OpenFile(loc, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	err = os.Truncate(latestname, 0)
	if err != nil {
		log.Println(err)
	}
	latestfile, err := os.OpenFile(latestname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	writeloc := io.MultiWriter(
		datefile,
		os.Stdout,
		latestfile,
	)
	if err != nil {
		log.Println(err)
	}
	Log, err = logger.New("Logger", 0, writeloc)
	if err != nil {
		log.Println(err)
	}
	Log.SetFormat(fmat)
	created = true
}

/* Closing the logger.
Due to the logging object not playing nice, I couldn't use it as a reciever*/
func Close(l *logger.Logger) {
	l.SetFormat("%{message}")
	const layout = "01-02-2006 15:04:05"
	// Place now in the string.
	t := time.Now()
	l.Infof("\n\nEnd of log for %v", t.Format(layout))
	l.Info(strings.Repeat("-", 80))
	l.Info("\n\n")
}
