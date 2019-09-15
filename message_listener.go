package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
	"github.com/leekchan/accounting"

	"github.com/bwmarrin/discordgo"
)

type bet struct {
	Description string    `firestore:"description"`
	Unit        string    `firestore:"unit"`
	Amount      float64   `firestore:"amount"`
	Odds        float64   `firestore:"odds"`
	Won         bool      `firestore:"won"`
	Resolved    bool      `firestore:"resolved"`
	User        string    `firestore:"user"`
	Timestamp   time.Time `firestore:"timestamp"`
}

var betMap map[int]*firestore.DocumentSnapshot
var deleted []int

func messageListener(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	response, err := understandMessage(message.Content, message.Author.Username)

	if err != nil {
		session.ChannelMessageSend(message.ChannelID, err.Error())
	}

	if response != "" {
		session.ChannelMessageSend(message.ChannelID, response)
	}
}

func understandMessage(content, user string) (string, error) {
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
		return parseBet(data, user)
	case "open":
		return openBets(user), nil
	case "delete":
		return deleteBets(data, user), nil
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
	case "delete":
		retStr.WriteString("Usage: !delete <ID as provided by !open>")
	default:
		retStr.WriteString("BetTrackerBot general help:\n")
		retStr.WriteString("`!help`: Shows this message! Use !help <command\\> to get more specific help for anything below\n")
		retStr.WriteString("`!bet`: Commits a bet to the database\n")
		retStr.WriteString("`!open`: Gets a list of bets you have open\n")
		retStr.WriteString("`!delete`: Deletes open bets\n")
		retStr.WriteString("`!ping`: Pong!")
	}
	return retStr.String()
}

func parseBet(data, user string) (string, error) {
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
		currencySymbol = 'u'
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

	err = placeBet(math.Floor(betAmount*100)/100, math.Floor(odds*100)/100, dataParts[2], string(currencySymbol), user)

	if err != nil {
		return "", errors.New("Couldn't store bet: " + err.Error())
	}

	if inUnits {
		return "Placed a bet of " + fmt.Sprintf("%.2f", betAmount) +
			" units @" + fmt.Sprintf("%.2f", odds) +
			" on " + dataParts[2], nil
	}
	ac := accounting.Accounting{Symbol: string(currencySymbol), Precision: 2}
	return "Placed a bet of " + ac.FormatMoney(betAmount) +
		" @" + fmt.Sprintf("%.2f", odds) +
		" on " + dataParts[2], nil
}

func isDigit(character rune) bool {
	value := character - '0'
	return value >= 0 && value <= 9
}

func placeBet(betAmount, odds float64, description, currencySymbol, user string) error {
	bets := client.Collection("Bets")
	_, _, err := bets.Add(context.Background(), bet{
		Description: description,
		Unit:        currencySymbol,
		Amount:      betAmount,
		Odds:        odds,
		User:        user,
		Won:         false,
		Resolved:    false,
		Timestamp:   time.Now(),
	})
	return err
}

func openBets(user string) string {
	bets := client.Collection("Bets")
	openBets := bets.Where("user", "==", user).Where("resolved", "==", false).OrderBy("timestamp", firestore.Asc).Documents(context.Background())

	i := 1
	betMap = make(map[int]*firestore.DocumentSnapshot)
	deleted = make([]int, 0)
	var retStr strings.Builder

	for {
		betDoc, err := openBets.Next()
		if err == iterator.Done {
			if i == 1 {
				return "No open bets found for " + user
			}
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			return "Error retrieving bets for " + user
		}

		var betData bet
		if err = betDoc.DataTo(&betData); err != nil {
			return "Couldn't parse a bet from the database"
		}

		var amount string
		if betData.Unit == "u" {
			amount = fmt.Sprintf("%.2f", betData.Amount) + betData.Unit
		} else {
			amount = betData.Unit + fmt.Sprintf("%.2f", betData.Amount)
		}
		retStr.WriteString(strconv.Itoa(i) + ": " + amount + " @" + fmt.Sprintf("%.2f", betData.Odds) + " on " + betData.Description + "\n")
		betMap[i] = betDoc
		i++
	}

	return retStr.String()
}

func deleteBets(data, user string) string {
	if len(betMap) == 0 {
		return "No open bets found to delete (you may need to !open first)"
	}

	betID, err := strconv.Atoi(data)
	if err != nil {
		return "Usage: !delete <ID as provided by !open>"
	}

	var betData bet
	if err = betMap[betID].DataTo(&betData); err != nil {
		return "Couldn't parse a bet from the database"
	}

	if user != betData.User {
		return "This isn't your bet to delete!"
	}

	for _, x := range deleted {
		if x == betID {
			return "Already deleted that bet"
		}
	}

	_, err = betMap[betID].Ref.Delete(context.Background())
	if err != nil {
		log.Printf(err.Error())
		return "Couldn't delete that bet from the database"
	}

	deleted = append(deleted, betID)
	return "Bet successfully deleted"
}
