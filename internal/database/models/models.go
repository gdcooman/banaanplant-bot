package models

import "gorm.io/gorm"

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
