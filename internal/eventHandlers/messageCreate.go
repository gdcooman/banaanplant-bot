package eventHandlers

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gdcooman/banaanplant-bot/internal/database/models"
	"gorm.io/gorm"
)

type MessageCreateHandler struct {
	db *gorm.DB
}

func NewMessageCreateHandler(db *gorm.DB) *MessageCreateHandler {
	return &MessageCreateHandler{
		db: db,
	}
}

func (h *MessageCreateHandler) Handler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	placeholders := map[string]string{
		"%author.mention%": message.Author.Mention(),
		"%author.name%":    message.Author.Username,
		"%bot.mention%":    session.State.User.Mention(),
		"%bot.name%":       session.State.User.Username,
	}

	var reaction models.CustomReaction
	result := h.db.Where("? ~* ('(\\W+|^)' || trigger || '(\\W+|$)')", message.Content).Order("length(trigger) desc").Limit(1).Find(&reaction)

	if result.RowsAffected == 1 {
		log.Println(fmt.Sprintf("Triggered reaction with ID: %d and trigger: %s", reaction.ID, reaction.Trigger))
		textReaction := reaction.TextReaction

		//Replace placeholders with actual values
		re, err := regexp.Compile("%[^%]*%")
		if err != nil {
			log.Fatal(err)
		}
		placeholderMatches := re.FindAllString(textReaction, -1)
		for _, placeholder := range placeholderMatches {
			textReaction = strings.Replace(textReaction, placeholder, placeholders[placeholder], 1)
		}

		//Add emoji reactions
		var emojiReactions []models.Emoji
		subQuery := h.db.Select("emoji_id").Where("custom_reaction_id = ?", reaction.ID).Table("reaction_emojis")
		h.db.Where("id in (?)", subQuery).Find(&emojiReactions)

		for _, emojiReac := range emojiReactions {
			session.MessageReactionAdd(message.ChannelID, message.ID, emojiReac.Unicode)
		}

		//Send reaction
		session.ChannelMessageSendReply(message.ChannelID, textReaction, message.MessageReference)
	}
}
