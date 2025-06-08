package plex

import "time"

type PlexPinResponse struct {
	ID               int       `json:"id"`
	Code             string    `json:"code"`
	ClientIdentifier string    `json:"clientIdentifier"`
	ExpiresAt        time.Time `json:"expiresAt"`
	AuthToken        string    `json:"authToken"`
}
