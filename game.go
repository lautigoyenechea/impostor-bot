package main

import (
	"fmt"
	"log"
	"maps"
	"math/rand"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
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

type Game struct {
	ID             string
	VoiceChannelID string
	Admin          Player
	Players        map[string]Player

	ImpostorID string
	Word       string

	VotingSession *VotingSession

	Ended bool

	TextChannelID string
}

func NewGame(voiceChannelID, textChannelID string, admin Player, players map[string]Player) *Game {
	return &Game{
		ID:             "#1",
		VoiceChannelID: voiceChannelID,
		Admin:          admin,
		Players:        players,
		ImpostorID:     pickImpostor(players),
		Word:           pickWord(),
		TextChannelID:  textChannelID,
	}
}

func (g *Game) IsAdmin(playerID string) bool {
	return g.Admin.ID == playerID
}

func (g *Game) IsImpostor(playerID string) bool {
	return g.ImpostorID == playerID
}

func (g *Game) End() {
	g.Ended = true
}

func (g *Game) SendWordToPlayers() {
	for id := range g.Players {
		if g.IsImpostor(id) {
			g.sendMessageToPlayer(id, fmt.Sprintf("Game %s - 🔪 You are the IMPOSTOR!", g.ID))
			continue
		}
		g.sendMessageToPlayer(id, fmt.Sprintf("Game %s - The word is: %s", g.ID, g.Word))
	}
}

func (g *Game) SendVotesToPlayers(guildID string) {
	votingSession := NewVotingSession(g, guildID)

	for id := range g.Players {
		g.sendComplexMessageToPlayer(id, &discordgo.MessageSend{
			Content:    "🗳️ **Voting Time!** Who do you think is the impostor?",
			Components: votingSession.Message.Rows,
		})
	}

	g.VotingSession = &votingSession
}

func (g *Game) AlivePlayersCount() int {
	return len(g.Players)
}

func (g *Game) AlivePlayersToText() string {
	var builder strings.Builder
	for _, p := range g.Players {
		builder.WriteString(fmt.Sprintf("- %s\n", p.Name))
	}
	return builder.String()
}

func (g *Game) EjectPlayer(playerID string) Player {
	player, _ := g.Players[playerID]
	delete(g.Players, playerID)
	return player
}

func (g *Game) sendMessageToPlayer(playerID, msg string) {
	dmChannel, err := session.UserChannelCreate(playerID)
	if err != nil {
		log.Printf("send message: creating player channel: id: %s\n", playerID)
		return
	}

	session.ChannelMessageSend(dmChannel.ID, msg)
}

func (g *Game) sendComplexMessageToPlayer(playerID string, msg *discordgo.MessageSend) {
	dmChannel, err := session.UserChannelCreate(playerID)
	if err != nil {
		log.Printf("send complex message: creating player channel: id: %s\n", playerID)
		return
	}

	session.ChannelMessageSendComplex(dmChannel.ID, msg)
}

func pickImpostor(players map[string]Player) string {
	playersIDs := slices.Collect(maps.Keys(players))
	return playersIDs[rand.Intn(len(playersIDs))]
}

func pickWord() string {
	return words[rand.Intn(len(words))]
}
