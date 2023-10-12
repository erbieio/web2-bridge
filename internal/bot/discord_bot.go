package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/erbieio/web2-bridge/config"
	"github.com/erbieio/web2-bridge/utils/discord"
	"github.com/erbieio/web2-bridge/utils/logger"
	"github.com/sirupsen/logrus"
)

type DiscordBot struct {
	Handler func(InputMessage) (OutputMessage, error)
}

func (bot *DiscordBot) App() string {
	return AppDiscord
}

func (bot *DiscordBot) Do() error {
	discord, err := discord.NewBot(config.GetDiscordConfig().BotToken, discordgo.IntentsGuildMessages, bot.MessageHandler)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("discord NewBot error")
		return err
	}
	err = discord.Open()
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("discord bot error")
		return err
	}
	return nil
}

func (bot *DiscordBot) MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.GuildID != "" && ((len(m.Mentions) == 0 || m.Mentions[0].ID != s.State.User.ID) || (m.ReferencedMessage != nil && m.ReferencedMessage.Author.ID != s.State.User.ID)) {
		return
	}
	logger.Logrus.Info(fmt.Sprintf("\nReceived message: %s from %s in channel %s\n", m.Content, m.Author.ID, m.ChannelID))
	compileRegex := regexp.MustCompile("transfernft\\s+([^\\s]*?)\\s+([^\\s]*)|mintnft")
	regArry := compileRegex.FindStringSubmatch(m.Content)
	if len(regArry) < 3 {
		return
	}

	msg, err := bot.Handler(InputMessage{
		App:       bot.App(),
		AuthorId:  fmt.Sprintf("%s/%s", bot.App(), m.Author.ID),
		MessageId: fmt.Sprintf("%s/%s", m.ChannelID, m.ID),
		Action:    regArry[0],
		Params:    regArry[1:],
	})
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("handle message error")
		return
	}
	msgInfo := strings.Split(msg.ReplyTo, "/")
	if len(msgInfo) < 2 {
		return
	}
	_, err = s.ChannelMessageSendReply(msgInfo[0], msg.Message, &discordgo.MessageReference{MessageID: msgInfo[1]})
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("discord bot send message error")
	}

}
