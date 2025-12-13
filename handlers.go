package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func StartCmdHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := session.State.Guild(i.GuildID)
	if err != nil {
		log.Fatalf("gettig guild: %v", err)
	}

	var voiceChannelID string

	for _, vs := range guild.VoiceStates {
		if vs.UserID == i.Member.User.ID {
			voiceChannelID = vs.ChannelID
			break
		}
	}

	if voiceChannelID == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ You must be in a voice channel to start a game!",
			},
		})
		return
	}

	var players []string
	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == voiceChannelID {
			players = append(players, vs.UserID)
		}
	}

	if len(players) < 3 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Must be at least 3 players in the server!",
			},
		})
		return
	}

	admin := Admin{
		ID:   i.Member.User.ID,
		Name: i.Member.User.Username,
	}

	game := NewGame(voiceChannelID, admin, players)
	server.AddGame(game)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Game %s started by %s", game.ID, game.Admin.Name),
		},
	})

	game.SendWordToPlayers()
}
