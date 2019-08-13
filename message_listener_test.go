package main

import "testing"

func TestUnderstandMessage(t *testing.T) {
	// Test empty message
	result, _ := understandMessage("")
	expectedResult := ""
	if result != expectedResult {
		t.Errorf("Empty message failed, expected %v, got %v", expectedResult, result)
	} else {
		t.Log("Empty message successfully tested")
	}

	// Test bad command
	_, err := understandMessage("!lolwhat")
	expectedResult = "Command lolwhat not found"
	if err == nil || err.Error() != expectedResult {
		t.Errorf("Bad command failed")
	} else {
		t.Log("Bad command successfully tested")
	}

	// Test !ping
	result, _ = understandMessage("!ping")
	expectedResult = "Pong!"
	if result != expectedResult {
		t.Errorf("!ping failed, expected %v, got %v", expectedResult, result)
	} else {
		t.Log("!ping successfully tested")
	}

	// Test empty !help
	result, _ = understandMessage("!help")
	expectedResult = "BetTrackerBot general help:\n" +
		"`!help`: Shows this message! Use !help <command\\> to get more specific help for anything below\n" +
		"`!bet`: Commits a bet to the database\n" +
		"`!ping`: Pong!"
	if result != expectedResult {
		t.Error("!help failed (you probably didn't update the test)")
	} else {
		t.Log("!help successfully tested")
	}

	// Test !help bet
	result, _ = understandMessage("!help bet")
	expectedResult = "Usage: !bet <amount> <decimal odds> <description of bet>\n" +
		"Amount can either be in units (eg 3.5u) or $/£/€ (eg $3.50)\n" +
		"Placing a bet will prompt the bot to check in with you 24h later to see if you won"
	if result != expectedResult {
		t.Error("!help bet failed")
	} else {
		t.Log("!help bet successfully tested")
	}

	// Test incorrect !bet strings
	expectedResult = "Usage: !bet <amount> <decimal odds> <description of bet>"

	_, err = understandMessage("!bet")
	if err == nil || err.Error() != expectedResult {
		t.Error("!bet empty string failed")
	} else {
		t.Log("!bet successfully tested")
	}

	_, err = understandMessage("!bet 3u")
	if err == nil || err.Error() != expectedResult {
		t.Error("!bet no description failed")
	} else {
		t.Log("!bet no description successfully tested")
	}

	expectedResult = "Please use an amount in units (eg 3.50u) or $/£/€ (eg $3.50)"

	_, err = understandMessage("!bet $3.50u 1.0 Weird Bet Dude")
	if err == nil || err.Error() != expectedResult {
		t.Error("!bet $3.50u failed")
	} else {
		t.Log("!bet $3.50u successfully tested")
	}

	_, err = understandMessage("!bet 3.50 2.0 Tree Fiddy Whats?")
	if err == nil || err.Error() != expectedResult {
		t.Error("!bet 3.50 failed")
	} else {
		t.Log("!bet 3.50 successfully tested")
	}

	expectedResult = "Couldn't parse bet amount"

	_, err = understandMessage("!bet $3.5s3g0 3.0 That's very odd")
	if err == nil || err.Error() != expectedResult {
		t.Error("!bet currency value not parsed failed")
	} else {
		t.Log("!bet currency value not parsed successfully tested")
	}

	_, err = understandMessage("!bet 3.5s3g0u 4.0 That's very odd, but in units")
	if err == nil || err.Error() != expectedResult {
		t.Error("!bet unit value not parsed failed")
	} else {
		t.Log("!bet unit value not parsed successfully tested")
	}

	_, err = understandMessage("!bet -3.50u 4.0 Don't bet negative money")
	expectedResult = "Bet amount negative"
	if err == nil || err.Error() != expectedResult {
		t.Error("!bet bet amount too low failed")
	} else {
		t.Log("!bet bet amount too low successfully tested")
	}

	_, err = understandMessage("!bet 3.5u 1.g0 Weird odds")
	expectedResult = "Couldn't parse odds"
	if err == nil || err.Error() != expectedResult {
		t.Error("!bet odds not parsed failed")
	} else {
		t.Log("!bet odds not parsed successfully tested")
	}

	_, err = understandMessage("!bet 3.5u 0.5 Too low")
	expectedResult = "Odds too low (<1.01)"
	if err == nil || err.Error() != expectedResult {
		t.Error("!bet odds too low failed")
	} else {
		t.Log("!bet odds too low successfully tested")
	}

	// !bet happy paths
	result, _ = understandMessage("!bet 3.5u 2.35 A good bet")
	expectedResult = "I would place a bet of 3.50 units @2.35 on A good bet, but this is too early in development"
	if result != expectedResult {
		t.Error("!bet units failed")
	} else {
		t.Log("!bet units successfully tested")
	}

	result, _ = understandMessage("!bet $3.50 2.35 Another great bet")
	expectedResult = "I would place a bet of $3.50 @2.35 on Another great bet, but this is too early in development"
	if result != expectedResult {
		t.Error("!bet money failed")
	} else {
		t.Log("!bet money successfully tested")
	}
}
