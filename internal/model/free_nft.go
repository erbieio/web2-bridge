package model

import "time"

const (
	TxStatusDefault = 0
	TxStatusSuccess = 1
	TxStatusFail    = 2
)

type FreeNft struct {
	Id             int       `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	App            string    `gorm:"column:app;NOT NULL"`                       // appname
	MintChannelId  string    `gorm:"column:mint_channel_id;NOT NULL"`           // mint消息channelid
	MintTxHash     string    `gorm:"column:mint_tx_hash;NOT NULL"`              // mint交易hash
	Creator        string    `gorm:"column:creator;NOT NULL"`                   // nft创建人标识
	TokenId        string    `gorm:"column:token_id;NOT NULL"`                  // nft唯一标识
	MintStatus     int       `gorm:"column:mint_status;default:0;NOT NULL"`     // 0-默认，1-成功，2-失败
	TransferTxHash string    `gorm:"column:transfer_tx_hash;NOT NULL"`          // transfer交易hash
	TransferStatus int       `gorm:"column:transfer_status;default:0;NOT NULL"` // 0-默认，1-成功，2-失败
	Created        time.Time `gorm:"column:created;default:CURRENT_TIMESTAMP;NOT NULL"`
	Updated        time.Time `gorm:"column:updated;default:CURRENT_TIMESTAMP;NOT NULL"`
}

func (m *FreeNft) TableName() string {
	return "free_nft"
}
