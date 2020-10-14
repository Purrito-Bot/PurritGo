package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	ng "github.com/djaustin/name-generator"
	"github.com/purrito-bot/purrigo/voice"
)

var generator ng.NameGenerator

var meowBuffer = make([][]byte, 0)

type parsedCommand struct {
	session *discordgo.Session
	message *discordgo.MessageCreate
	command string
	args    []string
}

func init() {
	// Set up Markov chains for name generation
	generator = ng.New()

	fileInfo, err := ioutil.ReadDir("names")
	if err != nil {
		log.Panicln("Unable to read names directory", err.Error())
	}

	for _, fi := range fileInfo {
		names := []string{}
		f, err := os.Open("names/" + fi.Name())
		if err != nil {
			log.Panicln("Unable to open", fi.Name(), err.Error())
		}
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&names)
		if err != nil {
			log.Panicln("Unable to decode JSON names file", err.Error())
		}
		generator.SeedData(strings.Split(fi.Name(), ".")[0], names)
	}

	meowBuffer, err = voice.LoadSound("meow.dca")
	if err != nil {
		log.Panicln("Cannot load sound", err.Error())
	}
}

func main() {
	token := flag.String("t", "", "Bot Token")
	flag.Parse()

	// Check a token was provided
	if *token == "" {
		*token = os.Getenv("DISCORD_TOKEN")
		if *token == "" {
			flag.Usage()
			return
		}
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
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates)

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

	// Ignore all messages created by the bot itself or without the command prefix
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, "go ") {
		return
	}

	command := parseCommand(s, m)

	switch command.command {
	case "name":
		handleName(command)
	case "meow":
		s.ChannelMessageSendTTS(m.ChannelID, "meow")
	case "mirror":
		s.ChannelMessageSend(m.ChannelID, m.Author.AvatarURL(""))
	case "show":
		handleShow(command)
	case "speak":
		handleSpeak(command)
	}
}

func parseCommand(s *discordgo.Session, m *discordgo.MessageCreate) parsedCommand {
	splitMessage := strings.Split(m.Content, " ")[1:]
	command := parsedCommand{
		session: s,
		message: m,
		command: splitMessage[0],
	}

	if len(splitMessage) > 1 {
		command.args = splitMessage[1:]
	}

	return command

}
