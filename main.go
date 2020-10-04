package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := flag.String("t", "", "Bot Token")
	flag.Parse()

	// Check a token was provided
	if *token == "" {
		flag.Usage()
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + *token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// We only need to receive messages at this point
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// Cleanly close down the Discord session when main() exits
	defer dg.Close()

	fmt.Println("Purrito is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}

// messageCreate is a handler called whenever a messages is received on a channel the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore all messages not using the command prefix
	if !strings.HasPrefix(m.Content, "go") {
		return
	}

	// Strip the prefix from he command
	command := strings.Split(m.Content, " ")[1]

	// If the message is "ping" reply with "Pong!"
	if command == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if command == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
