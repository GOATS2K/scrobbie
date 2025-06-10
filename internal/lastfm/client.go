package lastfm

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"sort"

	"github.com/goats2k/scrobbie/internal/config"
	"github.com/goats2k/scrobbie/internal/httpclient"
)

const API_URL string = "https://ws.audioscrobbler.com/2.0"

type LastFmClient struct {
	Config *config.LastFmConfig
}

func New(config *config.LastFmConfig) *LastFmClient {
	return &LastFmClient{
		Config: config,
	}
}

func addSignature(lc *LastFmClient, req *http.Request) {
	queryParams := req.URL.Query()

	var keys sort.StringSlice
	for key, _ := range queryParams {
		if key == "format" {
			continue
		}
		keys = append(keys, key)
	}
	keys.Sort()

	var signature string
	for _, key := range keys {
		signature += key
		signature += queryParams.Get(key)
	}
	signature = signature + lc.Config.ApiSecret

	signatureHash := md5.Sum([]byte(signature))
	signatureInHex := hex.EncodeToString(signatureHash[:])

	httpclient.AddQuery(req.URL, "api_sig", signatureInHex)
}

func doRequest[T any](lc *LastFmClient, req *http.Request) (*T, error) {
	httpclient.AddQuery(req.URL, "format", "json")
	httpclient.AddQuery(req.URL, "api_key", lc.Config.ApiKey)
	addSignature(lc, req)

	return httpclient.RunRequest[T](req)
}

func (lc *LastFmClient) GetRequestToken() (*LastFmToken, error) {
	req, _ := http.NewRequest(http.MethodGet, API_URL, nil)
	httpclient.AddQuery(req.URL, "method", "auth.gettoken")
	return doRequest[LastFmToken](lc, req)
}

func (lc *LastFmClient) GetSessionKey(token string) (*LastFmSession, error) {
	req, _ := http.NewRequest(http.MethodGet, API_URL, nil)
	httpclient.AddQuery(req.URL, "method", "auth.getSession")
	httpclient.AddQuery(req.URL, "token", token)
	return doRequest[LastFmSession](lc, req)
}
