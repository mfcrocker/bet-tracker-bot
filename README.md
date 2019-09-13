# Bet Tracker Bot
A discord bot for tracking bets!

## Usage
### Environment Variables
`DISCORD_TOKEN`: The Discord Token for the bot. Don't commit this!
`GOOGLE_PROJECT_ID`: The GCP Project ID your Firestore is in
`GOOGLE_APPLICATION_CREDENTIALS`: The path to your JSON containing a Service Account Key with R/W permissions to that Firestore

### Installing and Running
1. [Have a working Go installation](https://golang.org/doc/install)
2. `go get github.com/mfcrocker/bet-tracker-bot`
3. `go run github.com/mfcrocker/bet-tracker-bot`

## Commands
* `!help <command>`: Shows help for a particular command (or a list of commands if no command provided)
* `!bet <amount> <odds> <description of bet>`: Bets an amount in units (3.50u) or money ($3.50) at the specified odds on the specified bet, then stores that in the Firestore
* `!ping`: Pong!
