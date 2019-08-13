package main

import (
	"os"
	"testing"
)

func TestSetUpBot(t *testing.T) {
	// Happy path
	discord, err := setUpBot(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		t.Errorf("Couldn't set up Discord connection in test, error was %v", err)
	} else {
		t.Log("Discord connection happy path successfully tested")
	}

	discord.Close()

	// Sad path
	discord, err = setUpBot("whatrubbish")
	if err == nil {
		t.Error("Somehow you set up a Discord connection with a dodgy token")
	} else {
		t.Log("Discord connection sad path successfully tested")
	}
}
