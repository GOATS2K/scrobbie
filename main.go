package main

import (
	"errors"
	"io/fs"

	"github.com/charmbracelet/log"
	"github.com/goats2k/scrobbie/internal/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Info("failed to load .env", "error", err)
	}

	config := config.New()
	err := config.Read()
	if errors.Is(err, fs.ErrNotExist) {
		config.CreateConfigDirectory()
	}

	if config.Plex.NeedsConfiguring() {
		config.Plex.Configure()
	}

	if config.LastFm.NeedsConfiguring() {
		config.LastFm.Configure()
	}

	if err := config.Write(); err != nil {
		log.Fatal("failed to write config: %s", err)
	}

	// then, auth with lastfm

	// grab user libraries
	// ask user which library they're streaming music from
	// fetch all streams from said library
	// filter on last sync date (none by default? or allow users to set a sync range?)
}
