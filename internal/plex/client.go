package plex

import (
	"net/http"
	"net/url"

	"github.com/goats2k/scrobbie/internal/config"
)

type Plex interface {
	GetPin()
	GetUser(pin string)
	GetServers()
	GetLibraries()
	GetHomeUsers()
	GetPlaybackHistory()
}

type PlexClient struct {
	Config *config.PlexConfig
}

type Url url.URL

var PlexUrl string = "https://plex.tv"

func New(config *config.PlexConfig) *PlexClient {
	return &PlexClient{
		Config: config,
	}
}

func AddQuery(u *url.URL, key string, value string) {
	q := u.Query()
	q.Add(key, value)
	u.RawQuery = q.Encode()
}

func doRequest[T any](c *PlexClient, req http.Request) {
	if c.Config.AuthToken != "" {
		AddQuery(req.URL, "X-Plex-Token", c.Config.AuthToken)
	}
	AddQuery(req.URL, "X-Plex-Client-Identifier", c.Config.ClientIdentifier)

}

// Fetch a PIN to authenticate a user with via https://plex.tv/link
func (c *PlexClient) GetPin() {

}

// Fetches information for the user authenticated via [GetPin]
func (c *PlexClient) GetUser(pin string) {

}

func (c *PlexClient) GetServers() {

}

func (c *PlexClient) GetLibraries() {

}

func (c *PlexClient) GetUsers() {

}

func (c *PlexClient) GetPlaybackHistory() {
	// accountID is always 1 for the owner of the server
}
