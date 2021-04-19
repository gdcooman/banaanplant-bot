package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gdcooman/banaanplant-bot/internal/database"
	"github.com/gdcooman/banaanplant-bot/internal/eventHandlers"
	"github.com/gdcooman/banaanplant-bot/internal/instanceConfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// get secrets from instanceConfig.yml
	config, err := instanceConfig.ParseConfigFromYAMLFile("./configs/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	//Connect with postgres DB
	db, err := gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	database.Seed(db)

	//Create new Discord session
	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal(err)
		return
	}

	initSession(dg, db)

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

func initSession(dg *discordgo.Session, db *gorm.DB) {
	//Declare intents
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	//Add even handlers
	dg.AddHandler(eventHandlers.NewMessageCreateHandler(db).Handler)

}
