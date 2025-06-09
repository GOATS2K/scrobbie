package plex

import (
	"strings"
	"time"

	"github.com/goats2k/scrobbie/internal/util"
)

type PlexSessionHistoryResponse struct {
	MediaContainer MediaContainer `json:"MediaContainer,omitempty"`
}
type Metadata struct {
	Track            string        `json:"title,omitempty"`
	Album            string        `json:"parentTitle,omitempty"`
	Artist           string        `json:"grandparentTitle,omitempty"`
	Type             string        `json:"type,omitempty"`
	ViewedAt         UnixTimestamp `json:"viewedAt,omitempty"`
	AccountID        int           `json:"accountID,omitempty"`
	LibrarySectionID string        `json:"librarySectionID,omitempty"`
}
type MediaContainer struct {
	Size     int        `json:"size,omitempty"`
	Metadata []Metadata `json:"Metadata,omitempty"`
}

type UnixTimestamp time.Time

func (u *UnixTimestamp) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	unixTime, err := util.FromUnixTimestamp(s)
	if err != nil {
		return err
	}
	*u = UnixTimestamp(unixTime)
	return nil
}

func (u *UnixTimestamp) Time() time.Time {
	return time.Time(*u)
}
