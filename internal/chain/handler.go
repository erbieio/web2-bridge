package chain

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/exp/slices"

	"time"

	"github.com/erbieio/web2-bridge/config"
	"github.com/erbieio/web2-bridge/internal/bot"
	"github.com/erbieio/web2-bridge/internal/model"
	"github.com/erbieio/web2-bridge/utils/db/mysql"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var currentKey = time.Now().Format("2006-01-02")
var userMintCache = make(map[string]int64)

func mintContr(author string) (int64, error) {
	todayTime := time.Now()
	today := todayTime.Format("2006-01-02")
	if today > currentKey {
		currentKey = today
		userMintCache = make(map[string]int64)
	}
	nextDay := todayTime.AddDate(0, 0, 1).Format("2006-01-02")
	mintCount, ok := userMintCache[author]
	if !ok {
		var count int64
		err := mysql.GetDB().Model(&model.FreeNft{}).Where("creator = ? and mint_status in ? and created > ? and created < ?", author, []int{TxStatusSuccess, TxStatusDefault}, currentKey, nextDay).Count(&count).Error
		if err != nil {
			return 0, err
		}
		userMintCache[author] = count
	}
	return mintCount, nil

}

func MessageHandler(in bot.InputMessage) (bot.OutputMessage, error) {
	if in.IsMintNft() {
		//owner := common.BytesToAddress(crypto.Keccak256([]byte(in.AuthorId)))
		count, err := mintContr(in.AuthorId)
		if err != nil {
			return bot.OutputMessage{
				App:     in.App,
				ReplyTo: in.MessageId,
				Message: "Mint nft failed due to query error",
			}, err
		}
		if count >= int64(config.GetChainConfig().MaxMint) && !slices.Contains(config.GetChainConfig().WhiteList, in.AuthorId) {
			return bot.OutputMessage{
				App:     in.App,
				ReplyTo: in.MessageId,
				Message: "mint amount has exceed 10 limit",
			}, nil
		}
		tx, err := MintErbieNft(config.GetChainConfig().Rpc, config.GetChainConfig().NftAdminPriv, in.Params[0])
		if err != nil {
			return bot.OutputMessage{
				App:     in.App,
				ReplyTo: in.MessageId,
				Message: "Mint nft failed due to chain network error",
			}, err
		}
		messageInfo := strings.Split(in.MessageId, "/")
		channelId := ""
		if len(messageInfo) == 2 {
			channelId = messageInfo[0]
		}
		nft := model.FreeNft{
			App:           in.App,
			MintTxHash:    tx,
			Creator:       in.AuthorId,
			MintChannelId: channelId,
		}
		err = mysql.GetDB().Model(&model.FreeNft{}).Create(&nft).Error
		if err != nil {
			return bot.OutputMessage{
				App:     in.App,
				ReplyTo: in.MessageId,
				Message: "Mint nft failed",
			}, err
		}
		userMintCache[in.AuthorId] += 1
		mintTemplate := `✅ Success!📢 **Your mint transaction hash is:** 

➡️ %s`
		return bot.OutputMessage{
			App:     in.App,
			ReplyTo: in.MessageId,
			Message: fmt.Sprintf(mintTemplate, tx),
		}, nil

	} else if in.IsTransferNft() {
		//from := common.BytesToAddress(crypto.Keccak256([]byte(in.AuthorId)))
		//tokenId := in.Params[0]
		//to := common.HexToAddress(in.Params[1])
		nft := model.FreeNft{}
		err := mysql.GetDB().Where("token_id = ? and creator = ? and mint_status = ?", strings.ToLower(in.Params[0]), in.AuthorId, TxStatusSuccess).First(&nft).Error
		if err != nil {
			return bot.OutputMessage{
				App:     in.App,
				ReplyTo: in.MessageId,
				Message: "Nft not found",
			}, err
		}
		if nft.TransferStatus == TxStatusSuccess {
			return bot.OutputMessage{
				App:     in.App,
				ReplyTo: in.MessageId,
				Message: "Nft has been transferred",
			}, nil
		}
		if nft.TransferTxHash != "" && nft.TransferStatus == TxStatusDefault {
			return bot.OutputMessage{
				App:     in.App,
				ReplyTo: in.MessageId,
				Message: "Nft transfer is pending",
			}, nil
		}

		tx, err := TransferErbieNft(config.GetChainConfig().Rpc, config.GetChainConfig().NftAdminPriv, in.Params[0], in.Params[1])
		if err != nil {
			return bot.OutputMessage{
				App:     in.App,
				ReplyTo: in.MessageId,
				Message: "Transfer nft failed due to chain network error",
			}, err
		}
		err = mysql.GetDB().Model(&nft).Select("transfer_tx_hash", "transfer_status").Updates(model.FreeNft{TransferTxHash: tx, TransferStatus: TxStatusDefault}).Error
		if err != nil {
			return bot.OutputMessage{
				App:     in.App,
				ReplyTo: in.MessageId,
				Message: "Transfer nft failed",
			}, err
		}
		transferTemplate := `✅Success! The NFT transaction has been completed.

👇*Your transfer transaction hash is*:
%s
		`
		return bot.OutputMessage{
			App:     in.App,
			ReplyTo: in.MessageId,
			Message: fmt.Sprintf(transferTemplate, tx),
		}, nil

	}
	return bot.OutputMessage{
		App:     in.App,
		ReplyTo: in.MessageId,
		Message: "action not support yet",
	}, errors.New("action not support yet")
}

func MintErbieNft(nodeUrl, adminPriv, imageUrl string) (string, error) {
	ctx := context.Background()
	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		return "", err
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(adminPriv)
	if err != nil {
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("private key error.")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", err
	}
	gasLimit := uint64(100000)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}
	doubleGasPrice := (&big.Int{}).Mul(gasPrice, big.NewInt(2))
	transaction := MintNftTransaction{
		Type:    0,
		Royalty: 10,
		MetaURL: imageUrl,
		Version: "0.0.1",
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		return "", err
	}

	tx_data := append([]byte(ErbiePrefix), data...)
	tx := types.NewTransaction(nonce, fromAddress, big.NewInt(0), gasLimit, doubleGasPrice, tx_data)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

func TransferErbieNft(nodeUrl, adminPriv, tokenId, to string) (string, error) {
	ctx := context.Background()
	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		return "", err
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(adminPriv)
	if err != nil {
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("private key error.")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", err
	}
	gasLimit := uint64(100000)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}
	doubleGasPrice := (&big.Int{}).Mul(gasPrice, big.NewInt(2))
	transaction := TransferNftTransaction{
		Type:       1,
		NftAddress: tokenId,
		Version:    "0.0.1",
	}

	data, err := json.Marshal(transaction)
	if err != nil {
		return "", err
	}

	tx_data := append([]byte(ErbiePrefix), data...)
	toAddress := common.HexToAddress(to)
	tx := types.NewTransaction(nonce, toAddress, big.NewInt(0), gasLimit, doubleGasPrice, tx_data)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", err
	}
	return strings.ToLower(signedTx.Hash().String()), nil
}

func GetTxReceipt(nodeUrl string, txhash string) (*types.Receipt, error) {

	client, err := ethclient.Dial(nodeUrl)
	if err != nil {
		return nil, err
	}

	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(txhash))
	if err != nil {
		return nil, fmt.Errorf("TransactionReceipt %s %v", txhash, err)
	}

	if receipt == nil {
		return nil, fmt.Errorf("TransactionReceipt %s is null", txhash)
	}

	return receipt, nil
}

func GetAccountInfo(nodeUrl string, nftaddr common.Address, blockNumber *big.Int) (*Account, error) {
	client, err := rpc.Dial(nodeUrl)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	var result Account
	err = client.CallContext(context.Background(), &result, "eth_getAccountInfo", nftaddr, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return &result, err
}

func GetNftMetaUrl(nodeUrl string, nftaddr common.Address, blockNumber *big.Int) string {
	nft, err := GetAccountInfo(nodeUrl, nftaddr, blockNumber)
	if err != nil {
		return ""
	}
	if nft.Nft != nil {
		return nft.Nft.MetaURL
	}
	return ""
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	return hexutil.EncodeBig(number)
}
