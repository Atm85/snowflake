package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/atm85/snowflake/models"
	"github.com/bwmarrin/discordgo"
)

var (

	// Slice of thumbnails to cycle through.
	thumbnails = []string{
		"https://media.tenor.com/36vZDFMhH1AAAAAC/winter-snow.gif",
		"https://media.tenor.com/CxsIxqk24qUAAAAC/snow-day-winnie-the-pooh.gif",
		"https://media.tenor.com/8SbJtdBhM1sAAAAd/rex-snow.gif",
		"https://media.tenor.com/2VxUCA5PeuIAAAAC/snow-globe.gif",
		"https://media.tenor.com/hCjZ8X7XLtMAAAAC/snow-days-are-the-best-days-happy-snow-day.gif",
		"https://media.tenor.com/TsWtyC5QigMAAAAC/dog-snowbank.gif",
		"https://media.tenor.com/j1LBIhXrRNcAAAAC/fall-cute.gif",
		"https://media.tenor.com/0jur865wB9cAAAAS/snow-snowy.gif",
		"https://media.tenor.com/pkJvXAPVFkYAAAAC/cute-dinosaur-dino-winter.gif",
		"https://media.tenor.com/dMFMJnGg8h0AAAAC/letitsnow-frozen.gif",
		"https://media.tenor.com/c2zXVUQfCN8AAAAC/winter-peanuts.gif",
	}

	// Slice of commands to register and their options
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "snowfall",
			Description: "Shows the top snowfall.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "type",
					Description:  "The Forecast to show. day | season",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}

	// Map a Handler to the commands.
	handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"snowfall": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:

				title := ""
				description := ""
				data := i.ApplicationCommandData()
				rand.Seed(time.Now().Unix())

				switch data.Options[0].StringValue() {
				case "day":
					title = ":snowflake: Top Snowfall — Last 24 Hours :snowflake:"
					for _, result := range getResult().Snowfall {
						description += fmt.Sprintf("%s — `%s\"`\n\n", result.StationName, result.MaxSnowFall)
					}
				case "season":
					title = ":snowflake: Top Snow Depth Reports :snowflake:"
					for _, result := range getResult().Snowdepth {
						description += fmt.Sprintf("%s — `%s\"`\n\n", result.StationName, result.MaxSnowDepth)
					}
				}

				index := rand.Intn(len(thumbnails))
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							{
								Color: 0xffffff,
								Thumbnail: &discordgo.MessageEmbedThumbnail{
									URL:    thumbnails[index],
									Width:  512,
									Height: 512,
								},

								Footer: &discordgo.MessageEmbedFooter{
									Text:    "Powered by weatherusa.net",
									IconURL: "https://scontent-ord5-1.xx.fbcdn.net/v/t39.30808-6/310005176_540820237915943_2779069681157834873_n.png?_nc_cat=110&ccb=1-7&_nc_sid=09cbfe&_nc_ohc=Dxb0bP1j7tYAX_7WYsH&_nc_ht=scontent-ord5-1.xx&oh=00_AfBvYHQYOTOvH4EKr6dSkciL2zwW64Ezu7BrgD_jXuSFeQ&oe=63933165",
								},

								Title:       title,
								Description: description,
							},
						},
					},
				})

				if err != nil {
					panic(err)
				}

			case discordgo.InteractionApplicationCommandAutocomplete:
				data := i.ApplicationCommandData()
				choices := []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "day",
						Value: "day",
					},
					{
						Name:  "season",
						Value: "season",
					},
				}

				if data.Options[0].StringValue() != "" {
					choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
						Name:  data.Options[0].StringValue(),
						Value: "day",
					})
				}

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionApplicationCommandAutocompleteResult,
					Data: &discordgo.InteractionResponseData{
						Choices: choices,
					},
				})

				if err != nil {
					panic(err)
				}
			}
		},
	}
)

// Type containing bot spesific configuration.
type config struct {
	Token   string `json:"token"`
	GuildID string `json:"guild_id"`
}

// Reads from the config.json file and populates the 'config' struct.
func readConfig() (*config, error) {

	file, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	c := config{}
	err = json.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Gets the results of the top snowfall stats from the api.
func getResult() *models.Result {

	endpoint := "https://www.weatherusa.net/api/feed?type=winter_reports"

	req, _ := http.NewRequest("GET", endpoint, nil)
	res, _ := http.DefaultClient.Do(req)
	defer func() {
		res.Body.Close()
	}()

	var model *models.Result
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &model)
	return model
}

// Application entrypoint.
func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	config, err := readConfig()
	if err != nil {
		log.Fatalf("Unable to read config: %s\n", err)
	}

	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalf("Unable to initialize discord session: %s\n", err)
	}

	if err := session.Open(); err != nil {
		log.Fatalf("Unable to open connection to discord: %v\n", err)
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := handlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	commandRegistry, err := session.ApplicationCommandBulkOverwrite(session.State.User.ID, config.GuildID, commands)
	if err != nil {
		log.Fatalf("Unable to register commands: %v\n", err)
	}

	log.Println(session.State.User.Username + "#" + session.State.User.Discriminator + " is online!")

	<-stop
	log.Println("Gracefully shutting down...")
	for _, cmd := range commandRegistry {
		err := session.ApplicationCommandDelete(session.State.User.ID, config.GuildID, cmd.ID)
		if err != nil {
			log.Fatalf("Unable to delete %q command: %v", cmd.Name, err)
		}
	}
}
