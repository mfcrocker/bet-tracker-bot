package main

import (
	"errors"
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

	response, err := understandMessage(message.Content)

	if err != nil {
		session.ChannelMessageSend(message.ChannelID, err.Error())
	}

	if response != "" {
		session.ChannelMessageSend(message.ChannelID, response)
	}
}

func understandMessage(content string) (string, error) {
	if content == "" || !strings.HasPrefix(content, "!") {
		return "", nil
	}

	messageParts := strings.SplitN(content[len(content)-(len(content)-1):], " ", 2)
	command := messageParts[0]
	data := ""
	if len(messageParts) > 1 {
		data = messageParts[1]
	}

	switch command {
	case "help":
		return help(data), nil
	case "ping":
		return "Pong!", nil
	case "ding":
		return "8=====D", nil
	case "bet":
		return parseBet(data)
	default:
		return "", errors.New("Command " + command + " not found")
	}
}

func help(data string) string {
	var retStr strings.Builder
	switch data {
	case "bet":
		retStr.WriteString("Usage: !bet <amount> <decimal odds> <description of bet>\n")
		retStr.WriteString("Amount can either be in units (eg 3.5u) or $/£/€ (eg $3.50)\n")
		retStr.WriteString("Placing a bet will prompt the bot to check in with you 24h later to see if you won")
	default:
		retStr.WriteString("BetTrackerBot general help:\n")
		retStr.WriteString("`!help`: Shows this message! Use !help <command\\> to get more specific help for anything below\n")
		retStr.WriteString("`!bet`: Commits a bet to the database\n")
		retStr.WriteString("`!ping`: Pong!")
	}
	return retStr.String()
}

func parseBet(data string) (string, error) {
	// Expecting 3 data components: bet size, decimal odds and description
	dataParts := strings.SplitN(data, " ", 3)
	if (len(dataParts)) != 3 {
		return "", errors.New("Usage: !bet <amount> <decimal odds> <description of bet>")
	}

	inUnits := strings.HasSuffix(dataParts[0], "u")
	var currencySymbol = []rune(dataParts[0])[0]

	// No negative bet sizes
	if strings.HasPrefix(dataParts[0], "-") || strings.HasPrefix(dataParts[0][1:], "-") {
		return "", errors.New("Bet amount negative")
	}

	// You need 1 of a currrency symbol or units - not 0, not 2!
	if (inUnits && !isDigit(currencySymbol)) || (!inUnits && isDigit(currencySymbol)) {
		return "", errors.New("Please use an amount in units (eg 3.50u) or $/£/€ (eg $3.50)")
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
		return "", errors.New("Couldn't parse bet amount")
	}

	odds, err := strconv.ParseFloat(dataParts[1], 64)
	if err != nil {
		return "", errors.New("Couldn't parse odds")
	}

	if odds < 1.01 {
		return "", errors.New("Odds too low (<1.01)")
	}

	if inUnits {
		return "I would place a bet of " + fmt.Sprintf("%.2f", betAmount) +
			" units @" + fmt.Sprintf("%.2f", odds) +
			" on " + dataParts[2] +
			", but this is too early in development", nil
	}
	ac := accounting.Accounting{Symbol: string(currencySymbol), Precision: 2}
	return "I would place a bet of " + ac.FormatMoney(betAmount) +
		" @" + fmt.Sprintf("%.2f", odds) +
		" on " + dataParts[2] +
		", but this is too early in development", nil
}

func isDigit(character rune) bool {
	value := character - '0'
	return value >= 0 && value <= 9
}
