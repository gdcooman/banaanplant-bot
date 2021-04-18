package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gdcooman/banaanplant-bot/pkg/emoji"
	"gopkg.in/yaml.v3"
)

type instanceConfig struct {
	Token string `yaml:"Token"`
}

func getConfig() (c instanceConfig, err error) {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &c)
	return
}

func main() {
	// get secrets from config.yml
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	//Create new Discord session
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal(err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Banaanplant bot is now running! Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	if message.Content == "https://tenor.com/view/good-morning-vietnam-robin-williams-classic-announcer-radio-gif-4844905" {
		log.Println("Good morning Vietnam")
		session.MessageReactionAdd(message.ChannelID, message.ID, emoji.Red_heart)
		gm := fmt.Sprintf("Goeiemorgen, %s", message.Author.Mention())
		session.ChannelMessageSendReply(message.ChannelID, gm, message.MessageReference)
	}
}
