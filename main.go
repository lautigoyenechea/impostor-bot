package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	AUTH_TOKEN string
	GUILD_ID   string
)

var session *discordgo.Session

var server *Server

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	AUTH_TOKEN = os.Getenv("AUTH_TOKEN")
	GUILD_ID = os.Getenv("GUILD_ID")

	session, err = discordgo.New("Bot " + AUTH_TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	server = NewServer()
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "start",
			Description: "Start a new game",
		},
		{
			Name:        "vote",
			Description: "Start a voting session",
		},
	}
)

func main() {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("bot is running")
	})

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}

		case discordgo.InteractionMessageComponent:
			voteClickHandler(s, i)
		}
	})

	err := session.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	for _, cmd := range commands {
		_, err := session.ApplicationCommandCreate(session.State.User.ID, "", cmd)
		if err != nil {
			log.Fatalf("cannot create slash command %q: %v", cmd.Name, err)
		}
	}

	log.Println("Impostor Bot is online!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}
