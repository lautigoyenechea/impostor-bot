package main

import (
	"fmt"
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

	server = new(Server)
}

var (
	cmd = &discordgo.ApplicationCommand{
		Name:        "start",
		Description: "Start a new game",
	}
)

func main() {
	session.AddHandler(func(s *discordgo.Session, i *discordgo.Ready) {
		fmt.Println("Bot is ready!")
		_, err := session.ApplicationCommandCreate(session.State.User.ID, GUILD_ID, cmd)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Commands created!")
	})

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name == "start" {
			StartCmdHandler(s, i)
			return
		}
	})

	err := session.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	log.Println("Impostor Bot is online!")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}
