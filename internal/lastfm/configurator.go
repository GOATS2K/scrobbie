package lastfm

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/goats2k/scrobbie/internal/config"
)

type LastFmConfigurator struct {
	Config *config.LastFmConfig
}

func NewConfigurator(config *config.LastFmConfig) *LastFmConfigurator {
	return &LastFmConfigurator{
		Config: config,
	}
}

func (lc *LastFmConfigurator) NeedsConfiguring() bool {
	return lc.Config.ApiKey == "" || lc.Config.ApiSecret == "" || lc.Config.SessionKey == ""
}

func (lc *LastFmConfigurator) Configure() {
	fmt.Println("Authenticate with Last.fm")
	color.Magenta("-------------------------")

	fmt.Printf("Please apply for a last.fm application here: %s\n\n",
		color.MagentaString("https://www.last.fm/api/account/create"))

	// 1. Set API Key
	huh.NewInput().
		Title("Input your API key").
		Value(&lc.Config.ApiKey).
		Run()

	// 2. Set API Secret
	huh.NewInput().
		Title("Input your shared secret").
		Value(&lc.Config.ApiSecret).
		Run()

	// 3. Log in
	client := New(lc.Config)
	token, err := client.GetRequestToken()
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	authUrl := fmt.Sprintf("https://last.fm/api/auth/?api_key=%s&token=%s", lc.Config.ApiKey, token.Token)
	fmt.Printf("Nice! Your credentials worked. Now, please login to the following URL in your browser.")
	fmt.Println()
	color.Magenta(authUrl)
	fmt.Println()
	fmt.Println("scrobbie will store your session token once you've logged in.")

	sessionKeyChan := make(chan *LastFmSession, 1)

	go func() {
		for {
			resp, err := client.GetSessionKey(token.Token)
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			sessionKeyChan <- resp
		}
	}()

	loginResponse := <-sessionKeyChan
	lc.Config.SessionKey = loginResponse.Session.Key
	color.Green("Hi %s! You've successfully logged in!", loginResponse.Session.Name)
}
