package main

import (
	"errors"
	"fmt"
	"io/fs"

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

	// plexClient := plex.NewClient(config.Plex, config.LastSyncDate)
	// history, err := plexClient.GetPlaybackHistory()
	// if err != nil {
	// 	color.Red(fmt.Sprintf("Failed to get playback history from Plex: %s", err.Error()))
	// 	os.Exit(1)
	// }

	// sample := history.MediaContainer.Metadata[:10]
	// for _, item := range sample {
	// 	msg := fmt.Sprintf("[%s] %s - %s - %s", time.Time.Format(item.ViewedAt.Time(), time.RFC3339), item.Artist, item.Track, item.Album)
	// 	color.Magenta(msg)
	// }

	// filter on last sync date (none by default? or allow users to set a sync range?)
}
