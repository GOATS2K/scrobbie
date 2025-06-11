package lastfm

type LastFmScrobbleResponse struct {
	Scrobbles Scrobbles `json:"scrobbles,omitempty"`
}
type Artist struct {
	Corrected string `json:"corrected,omitempty"`
	Text      string `json:"#text,omitempty"`
}
type Album struct {
	Corrected string `json:"corrected,omitempty"`
	Text      string `json:"#text,omitempty"`
}
type Track struct {
	Corrected string `json:"corrected,omitempty"`
	Text      string `json:"#text,omitempty"`
}
type IgnoredMessage struct {
	Code string `json:"code,omitempty"`
	Text string `json:"#text,omitempty"`
}
type AlbumArtist struct {
	Corrected string `json:"corrected,omitempty"`
	Text      string `json:"#text,omitempty"`
}
type Scrobble struct {
	Artist         Artist         `json:"artist,omitempty"`
	Album          Album          `json:"album,omitempty"`
	Track          Track          `json:"track,omitempty"`
	IgnoredMessage IgnoredMessage `json:"ignoredMessage,omitempty"`
	AlbumArtist    AlbumArtist    `json:"albumArtist,omitempty"`
	Timestamp      string         `json:"timestamp,omitempty"`
}
type Attr struct {
	Ignored  int `json:"ignored,omitempty"`
	Accepted int `json:"accepted,omitempty"`
}
type Scrobbles struct {
	Scrobble Scrobble `json:"scrobble,omitempty"`
	Attr     Attr     `json:"@attr,omitempty"`
}
