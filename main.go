package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	ng "github.com/djaustin/name-generator"
	"github.com/purrito-bot/purrigo/voice"
)

var drowNames = []string{"Aamaneus", "Acostant", "Adehémar", "Aimeric", "Aimerics", "Aimerigatz", "Aimeriguet", "Aimes", "Alas", "Alazais", "Albarics", "Aldrics", "Alfans", "Alfonzenc", "Alias", "Aliazars", "Allard", "Amaldric", "Amaldrics", "Amalvis", "Amaneus", "Amerig", "Ancelmes", "Ancelmetz", "Anfos", "Araimfres", "Arbert", "Arguis", "Armans", "Arnaud", "Arnaudos", "Arnaut", "Arnautz", "Arsius", "Audegers", "Austor", "Azalbertz", "Azemar", "Baset", "Baudois", "Bausas", "Beranger", "Berengiers", "Bernart", "Bernat", "Bernatz", "Bertrans", "Borel", "Bovert", "Burcan", "Cadmar", "Chatbert", "Chinon", "Crespi", "Daire", "Dalmatz", "Danain", "Dragan", "Dragonetz", "Drogos", "Ebratz", "Elad", "Emeric", "Enricx", "Espanel", "Espas", "Estaci", "Esteve", "Estotz", "Exuperi", "Fanjaus", "Feris", "Ferrandos", "Filipot", "Focaut", "Foilan", "Folquets", "Fortaner", "Frezols", "Fricor", "Gaidon", "Gailhard", "Galters", "Garnier", "Gaston", "Gastos", "Gaucelis", "Gaudifer", "Gaudis", "Gautiers", "Gervais", "Gilabert", "Gilabertz", "Girauda", "Giraudetz", "Girauds", "Girautz", "Girvais", "Girvaitz", "Gobert", "Godafres", "Gontrand", "Gualhartz", "Gui", "Guigo", "Guilabert", "Guilabertz", "Guilelmes", "Guilhamos", "Guilhelmes", "Guilhelmet", "Guilhelms", "Guilhem", "Guilheumes", "Guion", "Guios", "Guiotz", "Guiraud", "Guiraudos", "Guiraut", "Guis", "Haylon", "Hugues", "Imbert", "Inard", "Isarn", "Isarts", "Isoartz", "Isodard", "Izarns", "Jacques", "Jaques", "Jaufre", "Jaufres", "Jean", "Joan", "Joans", "Johan", "Johans", "Jordas", "Joris", "Josselin", "Joudain", "Lambert", "Lamberts", "Lanval", "Lozoïc", "Lozoïs", "Lucatz", "Luzia", "Mamert", "MartisAlgais", "Michels", "Milon", "Milos", "Miquel", "Nicolas", "Otes", "Otz", "Patrice", "Peire", "Perrin", "Peyre", "Pons", "Quinault", "Raimon", "Rainaut", "Rainautz", "Rainers", "Rainiers", "Ramon", "Ramons", "Raolf", "Raolfs", "Raüli", "Reiambalts", "Remi", "Ricals", "Ricartz", "Richart", "Riquers", "Riton", "Roberts", "Robertz", "Rogers", "Rogerx", "Rogier", "Rostains", "Rostans", "Rotger", "Rotgiers", "Rotlans", "Sauson", "Savarics", "Segui", "Serin", "Sevin", "Sevis", "Sicard", "Sicart", "Simo", "Simos", "Sornehan", "Tecin", "Tezis", "Thibaud", "Thosa", "Tibal", "Tibaut", "Tibout", "Titbaut", "Uc", "Ucs", "Ug", "Ugos", "Ugs", "Ugues", "Valeray", "Vezias", "Xavier"}

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
	generator.SeedData("drow", drowNames)
	var err error
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
