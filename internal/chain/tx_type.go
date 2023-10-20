package chain

type MintNftTransaction struct {
	Type      uint8  `json:"type"`
	Royalty   uint32 `json:"royalty,omitempty"`
	MetaURL   string `json:"meta_url,omitempty"`
	Exchanger string `json:"exchanger,omitempty"`
	Version   string `json:"version"`
}

type TransferNftTransaction struct {
	Type       uint8  `json:"type"`
	NftAddress string `json:"nft_address"`
	Version    string `json:"version"`
}
