package bot

type InputMessage struct {
	App       string   `json:"app"`
	AuthorId  string   `json:"author_id"`
	MessageId string   `json:"message_id"`
	Action    string   `json:"action"`
	Params    []string `json:"params"`
}

type OutputMessage struct {
	App     string `json:"app"`
	Message string `json:"message"`
	ReplyTo string `json:"reply_to"`
}

func (in *InputMessage) IsMintNft() bool {
	return in.Action == ActionMintNft
}

func (in *InputMessage) IsTransferNft() bool {
	return in.Action == ActionTransferNft
}
