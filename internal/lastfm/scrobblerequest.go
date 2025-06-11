package lastfm

type LastFmScrobbleRequest struct {
	Method     string `json:"method"`
	Artist     string `json:"artist"`
	Track      string `json:"track"`
	Album      string `json:"album"`
	Timestamp  string `json:"timestamp"`
	SessionKey string `json:"sk"`
	ApiKey     string `json:"api_key"`
	Format     string `json:"format"`
}
