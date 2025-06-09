package plex

import (
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/fatih/color"
	"github.com/goats2k/scrobbie/internal/config"
	"github.com/goats2k/scrobbie/internal/util"
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
	return c.Config.UserAuthToken == "" || c.Config.ServerAuthToken == "" || c.Config.ServerUrl == "" || c.Config.ClientIdentifier == "" || c.Config.LibrarySectionID == 0
}

func (c *PlexConfigurator) Configure() {
	c.Config.ClientIdentifier = generateClientIdentifier()

	client := NewClient(c.Config, time.Time{})

	// 1. create plex login pin
	// 2. wait for user to login
	authToken, err := authenticateUser(client)
	if err != nil || authToken == "" {
		color.Red("Failed to get auth token.")
		os.Exit(1)
	}
	// user's access token - valid for server if they're the owner
	c.Config.UserAuthToken = authToken

	// 3. fetch server list
	server, address := selectServer(client)
	c.Config.ServerUrl = address.URI
	c.Config.ServerAuthToken = server.AccessToken

	// 4. select library
	selectedLibrary := selectLibrary(client)
	c.Config.LibrarySectionID = selectedLibrary

	color.Green("Plex is configured.")
}

func selectLibrary(client *PlexClient) int {
	libraries, err := client.GetLibraries()
	if err != nil {
		color.Red(fmt.Sprintf("Failed to get libraries: %s", err.Error()))
		os.Exit(1)
	}
	musicLibraries := util.Filter(libraries.MediaContainer.Directory, func(t *PlexLibrarySectionDirectory) bool {
		return t.Type == "artist"
	})

	if len(musicLibraries) == 0 {
		color.Red("No libraries found.")
		os.Exit(2)
	}

	var (
		options         []huh.Option[string]
		selectedLibrary string
	)

	for _, library := range musicLibraries {
		optionText := fmt.Sprintf("%s - %s", library.Key, library.Title)
		option := huh.NewOption(optionText, library.Key)
		options = append(options, option)
	}

	prompt := huh.NewSelect[string]().
		Title("Select a music library.").
		Options(options...).
		Value(&selectedLibrary)

	if err := prompt.Run(); err != nil {
		color.Red("Prompt cancelled.")
		os.Exit(3)
	}

	libraryId, _ := strconv.Atoi(selectedLibrary)

	return libraryId
}

func selectServer(client *PlexClient) (*PlexResourceDevice, *PlexResourceDeviceConnection) {
	resources, err := client.GetResources()
	if err != nil {
		color.Red(fmt.Sprintf("Failed to get resource list: %s", err.Error()))
		os.Exit(1)
	}

	servers := util.Filter(resources.Devices, func(t *PlexResourceDevice) bool {
		return t.Provides == "server" && t.Owned == "1"
	})

	if len(servers) == 0 {
		color.Red("No owned servers found on your account.")
		color.Red("Unfortunately, this means that Scrobbie can not get your playback history.")
		os.Exit(2)
	}

	var options []huh.Option[*PlexResourceDevice]
	for _, server := range servers {
		title := server.Name
		options = append(options, huh.NewOption(title, server))
	}

	var (
		selectedServer  *PlexResourceDevice
		selectedAddress *PlexResourceDeviceConnection
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[*PlexResourceDevice]().
				Title("Choose your Plex server").
				Options(options...).
				Value(&selectedServer),

			huh.NewSelect[*PlexResourceDeviceConnection]().
				Title("How would you like to connect to your server?").
				Description("Remote is recommended for dedicated servers.").
				OptionsFunc(func() []huh.Option[*PlexResourceDeviceConnection] {
					var options []huh.Option[*PlexResourceDeviceConnection]
					for _, connection := range selectedServer.Connection {
						var title string
						title = connection.Address
						if connection.Local == "1" {
							title += " (Local)"
						} else {
							title += " (Remote)"
						}
						options = append(options, huh.NewOption(title, &connection))
					}
					return options
				}, &selectedServer).
				Value(&selectedAddress),
		),
	)

	form.Run()

	return selectedServer, selectedAddress
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
