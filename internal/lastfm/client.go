package lastfm

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"sort"

	"github.com/goats2k/scrobbie/internal/config"
	"github.com/goats2k/scrobbie/internal/httpclient"
)

const API_URL string = "https://ws.audioscrobbler.com/2.0/"

type LastFmClient struct {
	Config *config.LastFmConfig
}

func New(config *config.LastFmConfig) *LastFmClient {
	return &LastFmClient{
		Config: config,
	}
}

func addSignature(lc *LastFmClient, req *http.Request) *http.Request {
	if req.Method == http.MethodGet {
		return createQueryParamSignature(lc, req)
	}
	return createBodySignature(lc, req)
}

func createBodySignature(lc *LastFmClient, req *http.Request) *http.Request {
	reqBody, _ := io.ReadAll(req.Body)
	params, _ := url.ParseQuery(string(reqBody))

	signature := createSignature(params)
	signatureInHex := getSignatureHex(lc, signature)
	params.Add("api_sig", signatureInHex)
	newReq, _ := http.NewRequest(req.Method, req.URL.String(), bytes.NewBufferString(string(params.Encode())))
	return newReq
}

func createSignature(params url.Values) string {
	var keys sort.StringSlice
	for key := range params {
		if key == "format" {
			continue
		}
		keys = append(keys, key)
	}
	keys.Sort()

	var signature string
	for _, key := range keys {
		signature += key
		signature += params.Get(key)
	}
	return signature
}

func createQueryParamSignature(lc *LastFmClient, req *http.Request) *http.Request {
	queryParams := req.URL.Query()
	signature := createSignature(queryParams)
	signatureInHex := getSignatureHex(lc, signature)
	httpclient.AddQuery(req.URL, "api_sig", signatureInHex)

	return req
}

func getSignatureHex(lc *LastFmClient, signature string) string {
	signature = signature + lc.Config.ApiSecret
	signatureHash := md5.Sum([]byte(signature))
	signatureInHex := hex.EncodeToString(signatureHash[:])
	return signatureInHex
}

func doRequest[T any](lc *LastFmClient, req *http.Request) (*T, error) {
	if req.Method == http.MethodGet {
		httpclient.AddQuery(req.URL, "format", "json")
		httpclient.AddQuery(req.URL, "api_key", lc.Config.ApiKey)
		if lc.Config.SessionKey != "" {
			httpclient.AddQuery(req.URL, "sk", lc.Config.SessionKey)
		}
	}

	newReq := addSignature(lc, req)

	return httpclient.RunRequest[T](newReq)
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

func getIgnoredMessage(code string) string {
	switch code {
	case "1":
		return "Artist was ignored"
	case "2":
		return "Track was ignored"
	case "3":
		return "Timestamp was too old"
	case "4":
		return "Timestamp was too new"
	case "5":
		return "Daily scrobble limit exceeded"
	default:
		return ""
	}
}

func (lc *LastFmClient) Scrobble(track *LastFmScrobbleRequest) (*LastFmScrobbleResponse, error) {
	track.Method = "track.scrobble"
	track.SessionKey = lc.Config.SessionKey
	track.ApiKey = lc.Config.ApiKey
	track.Format = "json"

	// surely there's a better way of doing this
	data := url.Values{}

	data.Set("method", track.Method)
	data.Set("artist", track.Artist)
	data.Set("track", track.Track)
	data.Set("album", track.Album)
	data.Set("timestamp", track.Timestamp)
	data.Set("sk", track.SessionKey)
	data.Set("api_key", track.ApiKey)
	data.Set("format", track.Format)

	req, _ := http.NewRequest(http.MethodPost, API_URL, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := doRequest[LastFmScrobbleResponse](lc, req)
	if err != nil {
		return resp, err
	}

	resp.Scrobbles.Scrobble.IgnoredMessage.Text = getIgnoredMessage(resp.Scrobbles.Scrobble.IgnoredMessage.Code)
	return resp, nil
}
