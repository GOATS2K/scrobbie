package plex

type PlexLibrarySectionDirectory struct {
	Key       string `json:"key"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Agent     string `json:"agent"`
	Language  string `json:"language"`
	Directory bool   `json:"directory"`
	Hidden    int    `json:"hidden"`
	Location  []struct {
		ID   int    `json:"id"`
		Path string `json:"path"`
	} `json:"Location"`
}

type PlexLibrarySectionResponse struct {
	MediaContainer struct {
		Size      int                           `json:"size"`
		Directory []PlexLibrarySectionDirectory `json:"Directory"`
	} `json:"MediaContainer"`
}
