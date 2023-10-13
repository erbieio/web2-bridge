package chain

import (
	"github.com/erbieio/web2-bridge/internal/bot"
)

func MessageHandler(in bot.InputMessage) (bot.OutputMessage, error) {
	if in.IsMintNft() {
		//owner := common.BytesToAddress(crypto.Keccak256([]byte(in.AuthorId)))

	} else if in.IsTransferNft() {
		//from := common.BytesToAddress(crypto.Keccak256([]byte(in.AuthorId)))
		//tokenId := in.Params[0]
		//to := common.HexToAddress(in.Params[1])

	} else {
		//not support yet

	}
	return bot.OutputMessage{
		App:     in.App,
		ReplyTo: in.MessageId,
		Message: "success",
	}, nil
}
