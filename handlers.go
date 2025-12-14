package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"start":       startCmdHandler,
	"vote":        voteCmdHandler,
	"vote-clicks": voteClickHandler,
}

func startCmdHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := session.State.Guild(i.GuildID)
	if err != nil {
		log.Fatalf("start cmd handler: getting guild: %v", err)
	}

	voiceChannelID := getChannelID(guild, i)

	if voiceChannelID == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ You must be in a voice channel to start a game!",
			},
		})
		return
	}

	var admin Player
	players := make(map[string]Player)

	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == voiceChannelID {
			id := vs.UserID
			user, _ := session.User(vs.UserID)
			player := NewPlayer(id, user.Username)
			players[id] = player

			if id == i.Member.User.ID {
				admin = player
			}
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

	game := NewGame(voiceChannelID, admin, players)
	server.AddGame(voiceChannelID, game)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Game %s started by %s", game.ID, admin.Name),
		},
	})

	game.SendWordToPlayers()
}

func voteCmdHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	playerID := i.Member.User.ID

	guild, err := session.State.Guild(i.GuildID)
	if err != nil {
		log.Fatalf("vote cmd handler: getting guild: %v", err)
	}

	channelID := getChannelID(guild, i)
	game := server.Game(channelID)
	if game == nil {
		log.Printf("vote cmd handler: game not found for channel id %s\n", channelID)
		return
	}

	if game.Ended {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ There is not an active game. /start to start a new one.",
			},
		})
		return
	}

	if !game.IsAdmin(playerID) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("❌ You are not the game Admin! Only %s can start a voting session.", game.Admin.Name),
			},
		})
		return
	}

	game.SendVotesToPlayers()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "⏳ Waiting for all players to vote...",
		},
	})
}

func voteClickHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	if !strings.HasPrefix(i.MessageComponentData().CustomID, "vote:") {
		return
	}

	votedForID := strings.TrimPrefix(i.MessageComponentData().CustomID, "vote:")
	voterID := i.Member.User.ID

	guild, err := session.State.Guild(i.GuildID)
	if err != nil {
		log.Fatalf("vote click handler: getting guild: %v", err)
	}

	channelID := getChannelID(guild, i)
	game := server.Game(channelID)
	if game == nil {
		log.Printf("vote click handler: game not found for channel id %s\n", channelID)
		return
	}

	if game.Ended {
		log.Printf("vote click handler: game for channel id %s has ended\n", channelID)
		return
	}

	if game.VotingSession == nil {
		log.Printf("vote click handler: voting session is null: channel id: %s: voter id: %s", channelID, voterID)
		return
	}

	if !game.VotingSession.IsActive {
		log.Printf("vote click handler: voting session is not active: channel id: %s: voter id: %s", channelID, voterID)
		return
	}

	game.VotingSession.SaveVote(voterID, votedForID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ Vote recorded!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if game.VotingSession.VotesCount() != game.AlivePlayersCount() {
		return
	}

	game.VotingSession.Close()

	mostVotedPlayerID := game.VotingSession.GetMostVoted()
	if game.IsImpostor(mostVotedPlayerID) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "IMPOSTOR Ejected! VICTORY! 🏆",
			},
		})
		game.End()
		return
	}

	game.EjectPlayer(mostVotedPlayerID)

	if game.AlivePlayersCount() == 2 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "😈 IMPOSTOR VICTORY 🏆",
			},
		})
		game.End()
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "The game continues. There is 1 impostor among us...",
		},
	})

	game.VotingSession.Close()
}

func getChannelID(guild *discordgo.Guild, i *discordgo.InteractionCreate) string {
	for _, vs := range guild.VoiceStates {
		if vs.UserID == i.Member.User.ID {
			return vs.ChannelID
		}
	}
	return ""
}
