package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	fmt.Println(token)

	discord, err := setUpBot(token)
	if err != nil {
		fmt.Println("Error while setting up bot", err)
	}

	fmt.Println("Bet bot is up and running! Ctrl+C will kill it.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}
