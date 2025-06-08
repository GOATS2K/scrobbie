package plex

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/goats2k/scrobbie/internal/config"
)

type Plex interface {
	GetPin() (response *PlexPinResponse, err error)
	GetUser(pinId int) (response *PlexPinResponse, err error)
	GetResources() (resp *PlexResourcesResponse, err error)
	GetLibraries()
	GetPlaybackHistory()
}

type PlexClient struct {
	Config *config.PlexConfig
}

var PlexUrl string = "https://plex.tv"

func NewClient(config *config.PlexConfig) *PlexClient {
	return &PlexClient{
		Config: config,
	}
}

var _ *PlexClient = &PlexClient{}
var httpClient *http.Client = &http.Client{
	Jar:       http.DefaultClient.Jar,
	Transport: http.DefaultTransport,
	Timeout:   30 * time.Second,
}

func AddQuery(u *url.URL, key string, value string) {
	q := u.Query()
	q.Add(key, value)
	u.RawQuery = q.Encode()
}

func doRequest[T any](c *PlexClient, req *http.Request) (*T, error) {
	if c.Config.AuthToken != "" {
		AddQuery(req.URL, "X-Plex-Token", c.Config.AuthToken)
	}

	req.Header.Add("Accept", "application/json")
	AddQuery(req.URL, "X-Plex-Client-Identifier", c.Config.ClientIdentifier)

	for {
		resp, err := httpClient.Do(req)
		if resp.StatusCode == 429 {
			log.Warn("HTTP 429 hit! Sleeping for 10 seconds before retrying.", "url", req.URL)
			time.Sleep(10 * time.Second)
			continue
		}

		if err != nil || resp.StatusCode >= 400 && resp.StatusCode <= 599 && resp.StatusCode != 429 {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}

		defer resp.Body.Close()

		var response T
		if strings.Contains(resp.Header.Get("Content-Type"), "application/xml") {
			if err := xml.NewDecoder(resp.Body).Decode(&response); err != nil {
				return nil, err
			}
			return &response, nil
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}
		return &response, nil
	}
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
	AddQuery(req.URL, "includeHttps", "1")
	return doRequest[PlexResourcesResponse](c, req)
}

func (c *PlexClient) GetLibraries() {

}

func (c *PlexClient) GetUsers() {

}

func (c *PlexClient) GetPlaybackHistory() {
	// accountID is always 1 for the owner of the server
}
