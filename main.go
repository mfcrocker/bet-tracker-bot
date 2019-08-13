package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if len(token) == 0 {
		log.Fatal("Must set DISCORD_TOKEN env var")
	}
	fmt.Println(token)

	discord, err := setUpBot(token)
	if err != nil {
		log.Fatal("Error while setting up bot", err)
	}
	defer discord.Close()

	fmt.Println("Bet bot is up and running! Ctrl+C will kill it.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	select {
	case <-sc:
		if discord != nil {
			discord.Close()
		}
	}
}

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
