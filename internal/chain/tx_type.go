package chain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

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

type AccountNFT struct {
	//Account
	Name   string
	Symbol string
	//Price                 *big.Int
	//Direction             uint8 // 0:no_tx,1:by,2:sell
	Owner                 common.Address
	SNFTRecipient         common.Address
	NFTApproveAddressList common.Address
	//Auctions map[string][]common.Address
	// MergeLevel is the level of NFT merged
	MergeLevel  uint8
	MergeNumber uint32
	//PledgedFlag           bool
	//NFTPledgedBlockNumber *big.Int

	Creator   common.Address
	Royalty   uint16
	Exchanger common.Address
	MetaURL   string
}

type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash // merkle root of the storage trie
	CodeHash []byte
	Nft      *AccountNFT `rlp:"nil"`
	Extra    []byte
}
