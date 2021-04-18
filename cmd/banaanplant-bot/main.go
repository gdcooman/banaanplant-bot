package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Config.yml
type instanceConfig struct {
	Token string `yaml:"Token"`
	Dsn   string `yaml:"Dsn"`
}

//DB Models
type CustomReaction struct {
	gorm.Model
	Trigger        string
	TextReaction   string
	EmojiReactions []Emoji `gorm:"many2many:reaction_emojis;"`
}

type Emoji struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	Unicode string
}

//Global variables
var db *gorm.DB

func getConfig() (c instanceConfig, err error) {
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(data, &c)
	return
}

func main() {
	var err error
	// get secrets from config.yml
	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	//Connect with postgres DB
	db, err = gorm.Open(postgres.Open(config.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	//Migrate and seed DB
	//db.AutoMigrate(&CustomReaction{}, &Emoji{})
	//
	//redHeart := Emoji{Name: "Red Heart", Unicode: "\U00002764"}
	//kiNon := Emoji{Name: "KiNon", Unicode: "<:kiNon:810648437904769117>"}
	//pileOfPoo := Emoji{Name: "Pile of Poo", Unicode: "\U0001F4A9"}
	//dragon := Emoji{Name: "Dragon", Unicode: "\U0001F409"}
	//t_rex := Emoji{Name: "T-Rex", Unicode: "\U0001F996"}
	//snake := Emoji{Name: "Snake", Unicode: "\U0001F40D"}
	//partyPopper := Emoji{Name: "Party Popper", Unicode: "\U0001F389"}
	//partyingFace := Emoji{Name: "Partying Face", Unicode: "\U0001F973"}
	//
	//db.Create(&redHeart)
	//db.Create(&kiNon)
	//db.Create(&pileOfPoo)
	//db.Create(&dragon)
	//db.Create(&t_rex)
	//db.Create(&snake)
	//db.Create(&partyPopper)
	//db.Create(&partyingFace)
	//
	//db.Create(&CustomReaction{Trigger: "https://tenor.com/view/good-morning-vietnam-robin-williams-classic-announcer-radio-gif-4844905", TextReaction: "Goeiemorgen, %author.mention%", EmojiReactions: []Emoji{redHeart}})
	//db.Create(&CustomReaction{Trigger: "no", TextReaction: "no u", EmojiReactions: []Emoji{kiNon}})
	//db.Create(&CustomReaction{Trigger: "Merci <@!809761796344512514>", TextReaction: "Geiren gedoan %author.mention%", EmojiReactions: []Emoji{t_rex}})
	//db.Create(&CustomReaction{Trigger: "Hello there", TextReaction: "https://tenor.com/view/hello-there-general-kenobi-star-wars-grevious-gif-17774326", EmojiReactions: []Emoji{dragon}})
	//db.Create(&CustomReaction{Trigger: "Ayy", TextReaction: "lmao", EmojiReactions: []Emoji{snake}})
	//db.Create(&CustomReaction{Trigger: "Mathias", TextReaction: "mama bankkaart", EmojiReactions: []Emoji{pileOfPoo}})
	//db.Create(&CustomReaction{Trigger: "Tis vrijdag", TextReaction: "IT'S SATURDAY, SUNDAY WHAT??? \n https://www.youtube.com/watch?v=IHC1Ma6EBbE", EmojiReactions: []Emoji{partyPopper, partyingFace}})

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

	log.Println(message.Content)

	placeholders := map[string]string{
		"%author.mention%": message.Author.Mention(),
		"%author.name%":    message.Author.Username,
		"%bot.mention%":    session.State.User.Mention(),
		"%bot.name%":       session.State.User.Username,
	}

	searchStr := "\\m(" + strings.ReplaceAll(message.Content, " ", "|") + ")\\M"

	var possibleReactions []CustomReaction
	db.Where("trigger ~* ?", searchStr).Find(&possibleReactions)
	log.Println(possibleReactions)

	for _, reac := range possibleReactions {
		if strings.Contains(strings.ToLower(message.Content), strings.ToLower(reac.Trigger)) {
			log.Println(fmt.Sprintf("Triggerd reaction with ID: %d and trigger: %s", reac.ID, reac.Trigger))
			reaction := reac.TextReaction
			re, err := regexp.Compile("%[^%]*%")
			if err != nil {
				log.Fatal(err)
			}

			placeholderMatches := re.FindAllString(reaction, -1)
			for _, placeholder := range placeholderMatches {
				reaction = strings.Replace(reaction, placeholder, placeholders[placeholder], 1)
			}

			var emojiReactions []Emoji
			subQuery := db.Select("emoji_id").Where("custom_reaction_id = ?", reac.ID).Table("reaction_emojis")
			db.Where("id in (?)", subQuery).Find(&emojiReactions)

			for _, emojiReac := range emojiReactions {
				session.MessageReactionAdd(message.ChannelID, message.ID, emojiReac.Unicode)
			}

			session.ChannelMessageSendReply(message.ChannelID, reaction, message.MessageReference)
		}
	}
}
