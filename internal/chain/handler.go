package chain

import "github.com/erbieio/web2-bridge/internal/bot"

func MessageHandler(in bot.InputMessage) (bot.OutputMessage, error) {
	if in.IsMintNft() {

	} else if in.IsTransferNft() {

	} else {
		//not support yet

	}
	return bot.OutputMessage{
		App:     in.App,
		ReplyTo: in.MessageId,
		Message: "success",
	}, nil
}
