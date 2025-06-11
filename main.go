package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/goats2k/scrobbie/internal/config"
	"github.com/goats2k/scrobbie/internal/lastfm"
	"github.com/goats2k/scrobbie/internal/plex"
)

func main() {
	config := config.New()
	err := config.Read()
	if errors.Is(err, fs.ErrNotExist) {
		config.CreateConfigDirectory()
	}

	plexConfigurator := plex.NewConfigurator(config.Plex)
	if plexConfigurator.NeedsConfiguring() {
		plexConfigurator.Configure()
	}

	lastFmConfigurator := lastfm.NewConfigurator(config.LastFm)
	if lastFmConfigurator.NeedsConfiguring() {
		lastFmConfigurator.Configure()
	}

	if err := config.Write(); err != nil {
		color.Red(fmt.Sprintf("Failed to write config: %s", err))
	}

	plexClient := plex.NewClient(config.Plex, config.LastSyncDate)
	lastFmClient := lastfm.New(config.LastFm)

	history, err := plexClient.GetPlaybackHistory()
	if err != nil {
		color.Red(fmt.Sprintf("Failed to get playback history from Plex: %s", err.Error()))
		os.Exit(1)
	}

	sample := history.MediaContainer.Metadata
	for _, item := range sample {
		msg := fmt.Sprintf("[%s] %s - %s - %s", time.Time.Format(item.ViewedAt.Time(), time.RFC3339), item.Artist, item.Track, item.Album)
		color.Magenta(msg)

		scrobbleRequest := &lastfm.LastFmScrobbleRequest{
			Artist:    item.Artist,
			Track:     item.Track,
			Album:     item.Album,
			Timestamp: strconv.FormatInt(item.ViewedAt.Time().Unix(), 10),
		}
		scrobble, err := lastFmClient.Scrobble(scrobbleRequest)
		if err != nil {
			color.Red("Failed to scrobble track: %s", err)
		}
		if scrobble.Scrobbles.Scrobble.IgnoredMessage.Code != "0" {
			color.Yellow(fmt.Sprintf("Scrobble ignored: %s", scrobble.Scrobbles.Scrobble.IgnoredMessage.Text))
		} else {
			color.Green("Scrobble success!")
		}
		fmt.Println("-------------")
		config.LastSyncDate = item.ViewedAt.Time().Add(1 * time.Second)
		time.Sleep(1 * time.Second)
	}

	config.Write()

	// filter on last sync date (none by default? or allow users to set a sync range?)
}
