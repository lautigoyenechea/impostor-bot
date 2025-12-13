package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

var words = []string{
	"Mate",
	"Notebook",
	"Mesa",
	"Silla",
	"Arbol",
	"Oso",
	"Tigre",
	"Calle",
	"Moto",
}

type Admin struct {
	ID   string
	Name string
}

type Game struct {
	ID             string
	VoiceChannelID string
	Admin          Admin
	Players        []string
	StartedAd      time.Time

	ImpostorID string
	Word       string

	//mu sync.RWMutex
}

func NewGame(voiceChannelID string, admin Admin, players []string) *Game {
	return &Game{
		ID:             "#1",
		VoiceChannelID: voiceChannelID,
		Admin:          admin,
		Players:        players,
		StartedAd:      time.Now(),
		ImpostorID:     pickImpostor(players),
		Word:           pickWord(),
	}
}

func (g *Game) IsAdmin(playerID string) bool {
	return g.Admin.ID == playerID
}

func (g *Game) IsImpostor(playerID string) bool {
	return g.ImpostorID == playerID
}

func (g *Game) SendWordToPlayers() {
	for _, playerID := range g.Players {
		if g.IsImpostor(playerID) {
			g.sendMessageToPlayer(playerID, fmt.Sprintf("Game %s - 🔪 You are the IMPOSTOR!", g.ID))
			continue
		}
		g.sendMessageToPlayer(playerID, fmt.Sprintf("Game %s - The word is: %s", g.ID, g.Word))
	}
}

func (g *Game) sendMessageToPlayer(playerID, msg string) {
	dmChannel, err := session.UserChannelCreate(playerID)
	if err != nil {
		log.Printf("creating player channel: id: %s\n", playerID)
		return
	}

	session.ChannelMessageSend(dmChannel.ID, msg)
}

func pickImpostor(players []string) string {
	return players[rand.Intn(len(players))]
}

func pickWord() string {
	return words[rand.Intn(len(words))]
}
