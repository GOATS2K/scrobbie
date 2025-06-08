package config

import (
	"crypto/rand"
	"fmt"
)

func (c *PlexConfig) NeedsConfiguring() bool {
	return c.AuthToken == "" || c.ServerUrl == "" || c.ClientIdentifier == ""
}

func (c *PlexConfig) Configure() {
	c.ClientIdentifier = generateClientIdentifier()

	// 1. create plex login pin
	// 2. wait for user to login
	// 3. fetch server list
	// 4. save target server for history
}

func generateClientIdentifier() string {
	randText := rand.Text()
	identifier := fmt.Sprintf("scrobbie-%s", randText[:8])
	return identifier
}

func authenticateUser() {
	// Get PIN
	// Prompt user with pin
	// Wait for user to login
	// Set auth token
}

func getServers() {
	// get server list
}

func selectServer() {

}

func getLibraries() {

}

func selectLibrary() {

}
