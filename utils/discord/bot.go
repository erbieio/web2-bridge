package discord

import (
	"github.com/bwmarrin/discordgo"
)

func NewBot(token string, intents discordgo.Intent, handler ...interface{}) (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	discord.Identify.Intents = discordgo.MakeIntent(intents)
	for i := 0; i < len(handler); i++ {
		discord.AddHandler(handler[i])
	}
	return discord, nil
}
