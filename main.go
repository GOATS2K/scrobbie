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

	day := 24 * time.Hour
	twoWeeksAgo := time.Now().Add(-(14 * day))

	if config.LastSyncDate.IsZero() || config.LastSyncDate.Before(twoWeeksAgo) {
		fmt.Println("Last.fm only supports tracks that were played 2 weeks ago or earlier.")
		fmt.Println("https://support.last.fm/t/retroactively-scrobble-past-tracks-with-original-date/4588/36")
		fmt.Println()
		fmt.Printf("Getting tracks from: %s\n", twoWeeksAgo.Format(time.DateOnly))
		config.LastSyncDate = twoWeeksAgo
	}

	plexClient := plex.NewClient(config.Plex, config.LastSyncDate)
	lastFmClient := lastfm.New(config.LastFm)

	fmt.Println("Getting playback history from Plex...")
	history, err := plexClient.GetPlaybackHistory()
	if err != nil {
		color.Red(fmt.Sprintf("Failed to get playback history from Plex: %s", err.Error()))
		os.Exit(1)
	}

	if len(history.MediaContainer.Metadata) == 0 {
		color.Yellow("No new tracks found.")
		return
	}

	for _, item := range history.MediaContainer.Metadata {
		msg := fmt.Sprintf("[%s] %s - %s", time.Time.Format(item.ViewedAt.Time(), time.RFC3339), item.Artist, item.Track)
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
}
