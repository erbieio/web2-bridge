package bot

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/erbieio/web2-bridge/config"
	"github.com/erbieio/web2-bridge/internal/model"
	"github.com/erbieio/web2-bridge/utils/comfyui"
	"github.com/erbieio/web2-bridge/utils/db/mysql"
	"github.com/erbieio/web2-bridge/utils/discord"
	"github.com/erbieio/web2-bridge/utils/ipfs"
	"github.com/erbieio/web2-bridge/utils/logger"
	"github.com/richinsley/comfy2go/client"
	"github.com/sirupsen/logrus"
)

type DiscordBot struct {
	Handler func(InputMessage) (OutputMessage, error)
	Comfyui *client.ComfyClient
}

func (bot *DiscordBot) App() string {
	return AppDiscord
}

func (bot *DiscordBot) Do() error {
	discord, err := discord.NewBot(config.GetDiscordConfig().BotToken, discordgo.IntentsGuildMessages|discordgo.IntentDirectMessages, bot.CommandHandler)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("discord NewBot error")
		return err
	}

	err = discord.Open()
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("discord bot error")
		return err
	}
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "mint_nft",
			Description: "mint nft on erbie chain",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompts",
					Description: "image scenarioal description",
					Required:    true,
				},
			},
		},
		{
			Name:        "transfer_nft",
			Description: "transfer your nft",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "token-id",
					Description: "nft id",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "to",
					Description: "to address",
					Required:    true,
				},
			},
		},
		{
			Name:        "owned_nft",
			Description: "list your nft",
		},
	}
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", v)
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("create discord command error")
		}
		registeredCommands[i] = cmd
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
		AuthorId:  fmt.Sprintf("%s::%s", bot.App(), m.Author.ID),
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

func (bot *DiscordBot) CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"mint_nft": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			go func() {
				var imageCid string
				failedTemple := "Your NFT minting failed, please initiate the command again.<@%s>"
				authorId := ""
				if i.Member != nil {
					authorId = i.Member.User.ID
				} else if i.User != nil {
					authorId = i.User.ID
				} else {
					logger.Logrus.WithFields(logrus.Fields{"discordBody": i}).Error("discord cannot get author id")
					return
				}
				discordResp, err := s.InteractionResponse(i.Interaction)
				if err != nil {
					logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("get discord interaction response error")
					return
				}
				if _, ok := optionMap["prompts"]; ok {
					descrip := optionMap["prompts"].Value.(string)
					imageBytes, _, err := comfyui.Prompts2Image(bot.Comfyui, descrip)
					if err != nil {
						logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("Prompts2Image error")
						_, err = s.ChannelMessageEdit(discordResp.ChannelID, discordResp.ID, fmt.Sprintf(failedTemple, authorId))
						if err != nil {
							logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("edit discord message error")
							return
						}
						return
					}
					ipfsClient := ipfs.NewClient(config.GetIpfsConfig().Api)
					imageCid, err = ipfsClient.Add(bytes.NewBuffer(imageBytes))
					if err != nil {
						logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("upload image to ipfs error")
						_, err = s.ChannelMessageEdit(discordResp.ChannelID, discordResp.ID, fmt.Sprintf(failedTemple, authorId))
						if err != nil {
							logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("edit discord message error")
							return
						}
						return
					}
				} else {
					return
				}

				//metaStruct.Image = config.GetIpfsConfig().HttpGateway + imageCid
				//metaStr, _ := json.Marshal(metaStruct)
				//cid, err := ipfsClient.Add(strings.NewReader(string(metaStr)))
				//if err != nil {
				//	logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("upload ipfs error")
				//	return
				//}

				outMsg, err := bot.Handler(InputMessage{
					App:       bot.App(),
					AuthorId:  authorId,
					MessageId: fmt.Sprintf("%s/%s", i.ChannelID, i.ID),
					Action:    ActionMintNft,
					Params:    []string{config.GetIpfsConfig().HttpGateway + imageCid},
				})
				if err != nil {
					logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("handle message error")
					_, err = s.ChannelMessageEdit(discordResp.ChannelID, discordResp.ID, fmt.Sprintf(failedTemple, authorId))
					if err != nil {
						logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("edit discord message error")
						return
					}
					return
				}
				_, err = s.ChannelMessageEdit(discordResp.ChannelID, discordResp.ID, outMsg.Message)
				if err != nil {
					logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("edit discord message error")
					return
				}

			}()

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "⏰Wait a second.⏳I'm processing your request.",
				},
			})
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("command handler error")
			}
		},
		"transfer_nft": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			tokenId := ""
			if option, ok := optionMap["token-id"]; ok {
				tokenId = option.StringValue()
			}
			to := ""
			if option, ok := optionMap["to"]; ok {
				to = option.StringValue()
			}
			authorId := ""
			if i.Member != nil {
				authorId = i.Member.User.ID
			} else if i.User != nil {
				authorId = i.User.ID
			} else {
				logger.Logrus.WithFields(logrus.Fields{"discordBody": i}).Error("discord cannot get author id")
				return
			}
			outMsg, err := bot.Handler(InputMessage{
				App:       bot.App(),
				AuthorId:  authorId,
				MessageId: fmt.Sprintf("%s/%s", i.ChannelID, i.ID),
				Action:    ActionTransferNft,
				Params:    []string{tokenId, to},
			})
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("handle message error")
			}
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: outMsg.Message,
				},
			})
		},
		"owned_nft": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			authorId := ""
			if i.Member != nil {
				authorId = i.Member.User.ID
			} else if i.User != nil {
				authorId = i.User.ID
			} else {
				logger.Logrus.WithFields(logrus.Fields{"discordBody": i}).Error("discord cannot get author id")
				return
			}
			nfts := make([]*model.FreeNft, 0)
			err := mysql.GetDB().Model(&model.FreeNft{}).Where("creator = ? and mint_status = ? and transfer_status != ?", authorId, model.TxStatusSuccess, model.TxStatusSuccess).Limit(1000).Find(&nfts).Error
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("db error")
				return
			}
			ownedTokenIds := make([]string, 0)
			for _, v := range nfts {
				ownedTokenIds = append(ownedTokenIds, v.TokenId)
			}
			tokenIdStr := strings.Join(ownedTokenIds, `

`)
			replyTemplate := `***🗃️ Here are your NFT list on the ErbieChain.***  <@%s>

%s
			
**📈 You can initiate transactions using /transfer_nft.**`
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(replyTemplate, authorId, tokenIdStr),
				},
			})

		},
	}
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i)
	}
}

func editMessage(s *discordgo.Session, channelId string, messageId string, message string) (*discordgo.Message, error) {
	return s.ChannelMessageEdit(channelId, messageId, message)
}
