package plex

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh/spinner"
	"github.com/fatih/color"
	"github.com/goats2k/scrobbie/internal/config"
)

type PlexConfigurator struct {
	Config *config.PlexConfig
}

func NewConfigurator(config *config.PlexConfig) *PlexConfigurator {
	return &PlexConfigurator{
		Config: config,
	}
}

func (c *PlexConfigurator) NeedsConfiguring() bool {
	return c.Config.AuthToken == "" || c.Config.ServerUrl == "" || c.Config.ClientIdentifier == ""
}

func (c *PlexConfigurator) Configure() {
	c.Config.ClientIdentifier = generateClientIdentifier()

	client := NewClient(c.Config)

	// 1. create plex login pin
	// 2. wait for user to login
	authToken, err := authenticateUser(client)
	if err != nil || authToken == "" {
		color.Red("Failed to get auth token.")
		os.Exit(1)
	}

	// 3. fetch server list
	// selectedServer := selectServer(client)

	// 4. save target server for history
	// c.Config.ServerUrl = selectedServer
	c.Config.AuthToken = authToken
}

func selectServer(client *PlexClient) string {
	panic("unimplemented")
}

func authenticateUser(client *PlexClient) (string, error) {
	pin, err := client.GetPin()
	if err != nil {
		color.Red("Failed to get PIN from Plex.")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Authenticate with Plex")
	color.Magenta("----------------------")
	fmt.Println()
	fmt.Println("Please visit https://plex.tv/link and enter the code:")
	fmt.Println()
	color.Magenta(pin.Code)
	fmt.Println()

	successChan := make(chan string, 1)

	pinAction := func() {
		for {
			pinResponse, err := client.GetUser(pin.ID)
			if err != nil {
				color.Yellow(err.Error())
				time.Sleep(5 * time.Second)
				continue
			}
			if pinResponse.AuthToken != "" {
				color.Green("Successfully logged in to Plex!")
				successChan <- pinResponse.AuthToken
				return
			}

			time.Sleep(5 * time.Second)
		}
	}

	if err := spinner.New().Title("Waiting for PIN entry...").Action(pinAction).Run(); err != nil {
		return "", err
	}

	return <-successChan, nil

}

func generateClientIdentifier() string {
	randText := rand.Text()
	identifier := fmt.Sprintf("scrobbie-%s", randText[:8])
	return identifier
}
