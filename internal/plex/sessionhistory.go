package plex

import (
	"github.com/goats2k/scrobbie/internal/util"
)

type PlexSessionHistoryResponse struct {
	MediaContainer MediaContainer `json:"MediaContainer,omitempty"`
}
type Metadata struct {
	Track            string             `json:"title,omitempty"`
	Album            string             `json:"parentTitle,omitempty"`
	Artist           string             `json:"grandparentTitle,omitempty"`
	Type             string             `json:"type,omitempty"`
	ViewedAt         util.UnixTimestamp `json:"viewedAt,omitempty"`
	AccountID        int                `json:"accountID,omitempty"`
	LibrarySectionID string             `json:"librarySectionID,omitempty"`
}
type MediaContainer struct {
	Size     int        `json:"size,omitempty"`
	Metadata []Metadata `json:"Metadata,omitempty"`
}
