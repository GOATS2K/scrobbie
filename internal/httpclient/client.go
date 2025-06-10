package httpclient

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
)

var httpClient *http.Client = &http.Client{
	Jar:       http.DefaultClient.Jar,
	Transport: http.DefaultTransport,
	Timeout:   30 * time.Second,
}

const USER_AGENT string = "scrobbie/0.1 (go-http-client/1.1)"

func AddQuery(u *url.URL, key string, value string) {
	q := u.Query()
	q.Add(key, value)
	u.RawQuery = q.Encode()
}

func RunRequest[T any](req *http.Request) (*T, error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Set("User-Agent", USER_AGENT)
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
