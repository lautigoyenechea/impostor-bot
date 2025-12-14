package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Player struct {
	ID   string
	Name string
}

func NewPlayer(id, name string) Player {
	return Player{ID: id, Name: name}
}

func (p Player) VoteButtom(guildID string) discordgo.MessageComponent {
	return discordgo.Button{
		Label:    p.Name,
		Style:    discordgo.PrimaryButton,
		CustomID: fmt.Sprintf("vote:%s:%s", guildID, p.ID),
	}
}
