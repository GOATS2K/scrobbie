package plex

import (
	"fmt"

	"net/http"
	"strconv"

	"time"

	"github.com/goats2k/scrobbie/internal/config"
	"github.com/goats2k/scrobbie/internal/httpclient"
)

type Plex interface {
	GetPin() (response *PlexPinResponse, err error)
	GetUser(pinId int) (response *PlexPinResponse, err error)
	GetResources() (resp *PlexResourcesResponse, err error)
	GetLibraries() (resp *PlexLibrarySectionResponse, err error)
	GetPlaybackHistory()
}

type PlexClient struct {
	Config       *config.PlexConfig
	LastSyncDate time.Time
}

var PlexUrl string = "https://plex.tv"

func NewClient(config *config.PlexConfig, lastSyncDate time.Time) *PlexClient {
	return &PlexClient{
		Config:       config,
		LastSyncDate: lastSyncDate,
	}
}

var _ *PlexClient = &PlexClient{}

func doRequest[T any](c *PlexClient, req *http.Request) (*T, error) {
	if c.Config.UserAuthToken != "" && req.URL.Hostname() == "plex.tv" {
		httpclient.AddQuery(req.URL, "X-Plex-Token", c.Config.UserAuthToken)
	} else if c.Config.ServerAuthToken != "" {
		httpclient.AddQuery(req.URL, "X-Plex-Token", c.Config.ServerAuthToken)
	}

	httpclient.AddQuery(req.URL, "X-Plex-Client-Identifier", c.Config.ClientIdentifier)

	return httpclient.RunRequest[T](req)
}

// Fetch a PIN to authenticate a user with via https://plex.tv/link
func (c *PlexClient) GetPin() (resp *PlexPinResponse, err error) {
	route := "/api/v2/pins"
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", PlexUrl, route), nil)
	return doRequest[PlexPinResponse](c, req)
}

// Fetches information for the user authenticated via [GetPin]
func (c *PlexClient) GetUser(pinId int) (resp *PlexPinResponse, err error) {
	route := "/api/v2/pins"
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s/%d", PlexUrl, route, pinId), nil)
	return doRequest[PlexPinResponse](c, req)
}

func (c *PlexClient) GetResources() (resp *PlexResourcesResponse, err error) {
	route := "/api/resources"
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", PlexUrl, route), nil)
	httpclient.AddQuery(req.URL, "includeHttps", "1")
	return doRequest[PlexResourcesResponse](c, req)
}

func (c *PlexClient) GetLibraries() (resp *PlexLibrarySectionResponse, err error) {
	route := "/library/sections"
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", c.Config.ServerUrl, route), nil)
	return doRequest[PlexLibrarySectionResponse](c, req)
}

func (c *PlexClient) GetPlaybackHistory() (resp *PlexSessionHistoryResponse, err error) {
	route := "/status/sessions/history/all"

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", c.Config.ServerUrl, route), nil)
	if !c.LastSyncDate.IsZero() {
		unixTime := c.LastSyncDate.Unix()
		httpclient.AddQuery(req.URL, "viewedAt>", strconv.Itoa(int(unixTime)))
	}
	httpclient.AddQuery(req.URL, "librarySectionID", strconv.Itoa(c.Config.LibrarySectionID))
	// accountID is always 1 for the owner of the server
	httpclient.AddQuery(req.URL, "accountID", "1")

	return doRequest[PlexSessionHistoryResponse](c, req)
}
