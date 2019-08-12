package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/leekchan/accounting"

	"github.com/bwmarrin/discordgo"
)

func messageListener(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	response := understandMessage(message.Content)

	if response != "" {
		session.ChannelMessageSend(message.ChannelID, response)
	}
}

func understandMessage(content string) string {
	if content == "" || content[:1] != "!" {
		return ""
	}

	messageParts := strings.SplitN(content[len(content)-(len(content)-1):], " ", 2)
	command := messageParts[0]
	var data string
	if len(messageParts) > 1 {
		data = messageParts[1]
	} else {
		data = ""
	}

	switch command {
	case "help":
		return help(data)
	case "ping":
		return "Pong!"
	case "ding":
		return "8=====D"
	case "bet":
		return parseBet(data)
	default:
		return ""
	}
}

func help(data string) string {
	switch data {
	case "bet":
		return "Usage: !bet <amount> <decimal odds> <description of bet>\n" +
			"Amount can either be in units (eg 3.5u) or $/£/€ (eg $3.50)\n" +
			"Placing a bet will prompt the bot to check in with you 24h later to see if you won"
	default:
		return "BetTrackerBot general help:\n" +
			"`!help`: Shows this message! Use !help <command\\> to get more specific help for anything below\n" +
			"`!bet`: Commits a bet to the database\n" +
			"`!ping`: Pong!"
	}
}

func parseBet(data string) string {
	// Expecting 3 data components: bet size, decimal odds and description
	dataParts := strings.SplitN(data, " ", 3)
	if (len(dataParts)) != 3 {
		return "Usage: !bet <amount> <decimal odds> <description of bet>"
	}

	var inUnits = false
	if dataParts[0][len(dataParts[0])-1:] == "u" {
		inUnits = true
	}
	var currencySymbol = []rune(dataParts[0])[0]

	// No negative bet sizes
	if []rune(dataParts[0])[0] == '-' || []rune(dataParts[0])[1] == '-' {
		return "Bet amount zero or negative"
	}

	// You need 1 of a currrency symbol or units - not 0, not 2!
	if (inUnits && !isDigit(currencySymbol)) || (!inUnits && isDigit(currencySymbol)) {
		return "Please use an amount in units (eg 3.50u) or $/£/€ (eg $3.50)"
	}
	var betAmountString string
	if inUnits {
		betAmountString = strings.TrimRight(dataParts[0], "u")
	} else {
		_, i := utf8.DecodeRuneInString(dataParts[0])
		betAmountString = dataParts[0][i:]
	}

	betAmount, err := strconv.ParseFloat(betAmountString, 64)
	if err != nil {
		return "Couldn't parse bet amount"
	}

	odds, err := strconv.ParseFloat(dataParts[1], 64)
	if err != nil {
		return "Couldn't parse odds"
	}

	if odds < 1.01 {
		return "Odds too low (<1.01)"
	}

	if inUnits {
		return "I would place a bet of " + fmt.Sprintf("%.2f", betAmount) +
			" units @" + fmt.Sprintf("%.2f", odds) +
			" on " + dataParts[2] +
			", but this is too early in development"
	}
	ac := accounting.Accounting{Symbol: string(currencySymbol), Precision: 2}
	return "I would place a bet of " + ac.FormatMoney(betAmount) +
		" @" + fmt.Sprintf("%.2f", odds) +
		" on " + dataParts[2] +
		", but this is too early in development"
}

func isDigit(character rune) bool {
	validDigits := [10]int32{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	for _, b := range validDigits {
		if character == b {
			return true
		}
	}
	return false
}
