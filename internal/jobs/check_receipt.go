package jobs

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/erbieio/web2-bridge/config"
	"github.com/erbieio/web2-bridge/internal/chain"
	"github.com/erbieio/web2-bridge/internal/model"
	"github.com/erbieio/web2-bridge/utils/db/mysql"
	"github.com/erbieio/web2-bridge/utils/discord"
	"github.com/erbieio/web2-bridge/utils/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
)

func HandleMintReceipt() {
	discord, _ := discord.NewBot(config.GetDiscordConfig().BotToken, discordgo.IntentsGuildMessages|discordgo.IntentDirectMessages)
	tick := time.Tick(time.Second * 2)
	for range tick {
		nfts := make([]*model.FreeNft, 0)
		err := mysql.GetDB().Model(&model.FreeNft{}).Where("mint_status = ?", chain.TxStatusDefault).Limit(BatchCount).Find(&nfts).Error
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("db error")
			continue
		}

		for _, v := range nfts {
			hash := v.MintTxHash
			receipt, err := chain.GetTxReceipt(config.GetChainConfig().Rpc, hash)
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("chain network error")
				continue
			}
			tokenId := getTokenId(receipt.Logs)
			if receipt.Status == 0 || tokenId == "" {
				err = mysql.GetDB().Model(v).Select("mint_status").Updates(model.FreeNft{MintStatus: chain.TxStatusFail}).Error
				if err != nil {
					logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("update mint status error")
					continue
				}
			} else {
				err = mysql.GetDB().Model(v).Select("mint_status", "token_id").Updates(model.FreeNft{MintStatus: chain.TxStatusSuccess, TokenId: strings.ToLower(tokenId)}).Error
				if err != nil {
					logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("update mint status error")
					continue
				}
				go replyDiscordMessage(discord, tokenId, v.MintChannelId, v.Creator)
			}

		}
	}
}

func replyDiscordMessage(discord *discordgo.Session, tokenId string, channelId string, creatorId string) {
	metaUrl := chain.GetNftMetaUrl(config.GetChainConfig().Rpc, common.HexToAddress(tokenId), nil)
	if metaUrl == "" {
		return
	}
	resp, err := http.Get(metaUrl)
	if err != nil {
		logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("get nft meta url error")
		return
	}
	payloald := discordgo.MessageSend{}
	replyTemplate := `<@%s> ➡️Your NFT ID on the #ErbieChain is：%s

**Discord name of the🔎NFT creator: <@%s>**`
	payloald.Content = fmt.Sprintf(replyTemplate, creatorId, tokenId, creatorId)
	payloald.Files = []*discordgo.File{
		{
			Name:        "image.png",
			ContentType: "png",
			Reader:      resp.Body,
		},
	}
	discord.ChannelMessageSendComplex(channelId, &payloald)
}

func HandleTransferReceipt() {
	tick := time.Tick(time.Second * 2)
	for range tick {
		nfts := make([]*model.FreeNft, 0)
		err := mysql.GetDB().Model(&model.FreeNft{}).Where("transfer_status = ? and transfer_tx_hash != \"\"", chain.TxStatusDefault).Limit(BatchCount).Find(&nfts).Error
		if err != nil {
			logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("db error")
			continue
		}

		for _, v := range nfts {
			hash := v.TransferTxHash
			receipt, err := chain.GetTxReceipt(config.GetChainConfig().Rpc, hash)
			if err != nil {
				logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("chain network error")
				continue
			}
			if receipt.Status == 0 {
				err = mysql.GetDB().Model(v).Select("transfer_status").Updates(model.FreeNft{TransferStatus: chain.TxStatusFail}).Error
				if err != nil {
					logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("update transfer status error")
					continue
				}
			} else {
				err = mysql.GetDB().Model(v).Select("transfer_status").Updates(model.FreeNft{TransferStatus: chain.TxStatusSuccess}).Error
				if err != nil {
					logger.Logrus.WithFields(logrus.Fields{"Error": err}).Error("update transfer status error")
					continue
				}
			}

		}
	}
}

func getTokenId(logs []*types.Log) string {
	for _, log := range logs {
		if log.Removed {
			continue
		}
		if len(log.Topics) == 2 && log.Topics[0].String() == chain.MintNftTopic {
			return common.HexToAddress(log.Topics[1].String()).String()
		}

	}
	return ""
}
