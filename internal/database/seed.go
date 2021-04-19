package database

import (
	"github.com/gdcooman/banaanplant-bot/internal/database/models"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	db.Exec("drop table if exists reaction_emojis")
	db.Exec("drop table if exists custom_reactions")
	db.Exec("drop table if exists emojis")

	//Migrate and seed DB
	db.AutoMigrate(&models.CustomReaction{}, &models.Emoji{})

	redHeart := models.Emoji{Name: "Red Heart", Unicode: "\U00002764"}
	kiNon := models.Emoji{Name: "KiNon", Unicode: ":kiNon:810648437904769117"}
	pileOfPoo := models.Emoji{Name: "Pile of Poo", Unicode: "\U0001F4A9"}
	dragon := models.Emoji{Name: "Dragon", Unicode: "\U0001F409"}
	t_rex := models.Emoji{Name: "T-Rex", Unicode: "\U0001F996"}
	snake := models.Emoji{Name: "Snake", Unicode: "\U0001F40D"}
	partyPopper := models.Emoji{Name: "Party Popper", Unicode: "\U0001F389"}
	partyingFace := models.Emoji{Name: "Partying Face", Unicode: "\U0001F973"}
	darkPersonInMotorizedWheelchair := models.Emoji{Name: "Person in Motorized Wheelchair: Dark Skin Tone", Unicode: "\U0001F9D1\U0001F3FF\U0000200D\U0001F9BC"}

	db.Create(&redHeart)
	db.Create(&kiNon)
	db.Create(&pileOfPoo)
	db.Create(&dragon)
	db.Create(&t_rex)
	db.Create(&snake)
	db.Create(&partyPopper)
	db.Create(&partyingFace)
	db.Create(&darkPersonInMotorizedWheelchair)

	db.Create(&models.CustomReaction{Trigger: "https://tenor.com/view/good-morning-vietnam-robin-williams-classic-announcer-radio-gif-4844905", TextReaction: "Goeiemorgen, %author.mention%", EmojiReactions: []models.Emoji{redHeart}})
	db.Create(&models.CustomReaction{Trigger: "no", TextReaction: "no u", EmojiReactions: []models.Emoji{kiNon}})
	db.Create(&models.CustomReaction{Trigger: "Hello there", TextReaction: "https://tenor.com/view/hello-there-general-kenobi-star-wars-grevious-gif-17774326", EmojiReactions: []models.Emoji{dragon}})
	db.Create(&models.CustomReaction{Trigger: "Ayy", TextReaction: "lmao", EmojiReactions: []models.Emoji{snake}})
	db.Create(&models.CustomReaction{Trigger: "Mathias", TextReaction: "mama bankkaart", EmojiReactions: []models.Emoji{pileOfPoo}})
	db.Create(&models.CustomReaction{Trigger: "Tis vrijdag", TextReaction: "IT'S SATURDAY, SUNDAY WHAT??? \n https://www.youtube.com/watch?v=IHC1Ma6EBbE", EmojiReactions: []models.Emoji{partyPopper, partyingFace}})
	db.Create(&models.CustomReaction{Trigger: "Merci <@!809761796344512514>", TextReaction: "Geiren gedoan %author.mention%", EmojiReactions: []models.Emoji{t_rex}})
	db.Create(&models.CustomReaction{Trigger: "<@!809761796344512514>", TextReaction: "WUK IST \U0001F921", EmojiReactions: []models.Emoji{darkPersonInMotorizedWheelchair}})
	//Banaanplant Bot_DEV
	//db.Create(&models.CustomReaction{Trigger: "Merci <@!814516668389916672>", TextReaction: "Geiren gedoan %author.mention%", EmojiReactions: []models.Emoji{t_rex}})
	//db.Create(&models.CustomReaction{Trigger: "<@!814516668389916672>", TextReaction: "WUK IST \U0001F921", EmojiReactions: []models.Emoji{darkPersonInMotorizedWheelchair}})
}
