package main

import "github.com/bwmarrin/discordgo"

type Player struct {
	ID   string
	Name string
}

func NewPlayer(id, name string) Player {
	return Player{ID: id, Name: name}
}

func (p Player) VoteButtom() discordgo.MessageComponent {
	return discordgo.Button{
		Label:    p.Name,
		Style:    discordgo.PrimaryButton,
		CustomID: "vote-" + p.ID,
	}
}
