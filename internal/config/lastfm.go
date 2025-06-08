package config

func (c *LastFmConfig) NeedsConfiguring() bool {
	return c.SessionKey == "" || c.ApiKey == "" || c.ApiSecret == ""
}

func (c *LastFmConfig) Configure() {

}
