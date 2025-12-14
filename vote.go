package main

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type VotingSession struct {
	Message VoteMessage

	Votes    map[string]string
	IsActive bool

	mu sync.RWMutex
}

type VoteMessage struct {
	Rows []discordgo.MessageComponent
}

func NewVotingSession(game *Game, guildID string) VotingSession {
	if game == nil {
		log.Fatal("new voting session: game can not be nil")
		return VotingSession{}
	}

	return VotingSession{
		Message:  newVoteMessage(game, guildID),
		Votes:    make(map[string]string),
		IsActive: true,
	}
}

func (vs *VotingSession) SaveVote(voterID, voteForID string) {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	vs.Votes[voterID] = voteForID
}

func (vs *VotingSession) VotesCount() int {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	return len(vs.Votes)
}

func (vs *VotingSession) CountVotes() map[string]int {
	vs.mu.RLock()
	defer vs.mu.RUnlock()

	voteCounts := make(map[string]int)

	for _, votedForID := range vs.Votes {
		voteCounts[votedForID]++
	}

	return voteCounts
}

func (vs *VotingSession) GetMostVoted() string {
	voteCounts := vs.CountVotes()

	var mostVotedID string
	var maxVotes int

	for playerID, count := range voteCounts {
		if count > maxVotes {
			maxVotes = count
			mostVotedID = playerID
		}
	}

	return mostVotedID
}

func (vs *VotingSession) Close() {
	vs.IsActive = false
}

func newVoteMessage(game *Game, guildID string) VoteMessage {
	var buttoms []discordgo.MessageComponent

	for _, player := range game.Players {
		buttoms = append(buttoms, player.VoteButtom(guildID))
	}

	var rows []discordgo.MessageComponent
	lenButtoms := len(buttoms)

	for i := 0; i < lenButtoms; i += 5 {
		end := i + 5
		if end > lenButtoms {
			end = lenButtoms
		}

		row := discordgo.ActionsRow{
			Components: buttoms[i:end],
		}
		rows = append(rows, row)
	}

	return VoteMessage{Rows: rows}
}
