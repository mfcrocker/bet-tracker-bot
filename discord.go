package main

import (
	"github.com/bwmarrin/discordgo"
)

func setUpBot(token string) (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return discord, err
	}

	discord.AddHandler(messageListener)

	err = discord.Open()
	if err != nil {
		return discord, err
	}

	return discord, nil
}
