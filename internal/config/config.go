package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path"
	"time"

	"github.com/charmbracelet/log"
)

type LastFmConfig struct {
	ApiKey     string `json:"api_key"`
	ApiSecret  string `json:"api_secret"`
	SessionKey string `json:"session_key"`
}

type PlexConfig struct {
	ServerUrl        string `json:"server_url"`
	AuthToken        string `json:"auth_token"`
	LibrarySectionID int    `json:"library_section_id"`
	ClientIdentifier string `json:"client_identifier"`
}

type Config struct {
	LastFm       *LastFmConfig `json:"lastfm"`
	Plex         *PlexConfig   `json:"plex"`
	LastSyncDate time.Time     `json:"last_sync_date"`
}

func New() *Config {
	return &Config{
		Plex:         &PlexConfig{},
		LastFm:       &LastFmConfig{},
		LastSyncDate: time.Time{},
	}
}

func getConfigDir() string {
	var userConfigDir string

	_, err := os.Stat("/.dockerenv")
	if err == nil {
		log.Debug("docker container environment detected!")
		return "/config"
	}

	userConfigDir, err = os.UserConfigDir()
	if err == nil {
		return path.Join(userConfigDir, "scrobbie")
	}

	log.Warn("failed to get user config dir")
	userConfigDir, err = os.Getwd()
	if err != nil {
		log.Fatal("failed to get a config directory, aborting...", "error", err)
	}

	log.Debug("config path", "path", userConfigDir)

	return userConfigDir
}

func (c *Config) CreateConfigDirectory() (configDir string, err error) {
	configDirPath := getConfigDir()
	err = os.MkdirAll(configDirPath, 0755)
	if err == nil {
		log.Info("got config dir", "dir", configDirPath)
		return configDirPath, err
	}

	log.Error("failed to create config dir", "error", err)
	return "", err
}

func (c *Config) Write() error {
	configDirPath := getConfigDir()
	filePath := path.Join(configDirPath, "config.json")
	configFile, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer configFile.Close()

	enc := json.NewEncoder(configFile)
	enc.SetIndent("", "	")

	if err := enc.Encode(&c); err != nil {
		return err
	}

	log.Info("written config", "path", filePath)
	return nil
}

func (c *Config) Read() error {
	configDirPath := getConfigDir()
	filePath := path.Join(configDirPath, "config.json")
	_, err := os.Stat(filePath)

	if err != nil && errors.Is(err, fs.ErrNotExist) {
		return err
	}

	log.Info("loading config file...")
	configFile, err := os.Open(filePath)
	if err != nil {
		return err
	}

	json.NewDecoder(configFile).Decode(c)
	return nil
}
