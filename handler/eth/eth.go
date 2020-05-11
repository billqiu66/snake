package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lexkong/log"
	"math"
	"math/big"
)

func Balance(client *ethclient.Client, address string) (*big.Float, error) {
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Warnf("get balance err: %v", err)
		return nil, err
	}
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return ethValue, nil
}

func Transfer(client *ethclient.Client, private string, address string, num int64) (string, error) {
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Warnf("transfer err")
		return "", errors.New("transfer err")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}
	value := big.NewInt(num * 1000000000000000000)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}
	toAddress := common.HexToAddress(address)
	var data []byte
	// 未签名的交易
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}
	// 用私钥签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}
	// 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}
	txHash := signedTx.Hash().Hex()
	return txHash, nil
}

func checkBlockHeader(client *ethclient.Client, num *big.Int) (*types.Header, error) {
	header, err := client.HeaderByNumber(context.Background(), num)
	if err != nil {
		log.Warnf("check header err: %v", err)
		return nil, err
	}
	return header, nil
}
func checkBlock(client *ethclient.Client, num *big.Int) (*types.Block, error) {
	block, err := client.BlockByNumber(context.Background(), num)
	if err != nil {
		log.Warnf("check header err: %v", err)
		return nil, err
	}
	return block, nil
}

func checkTxCount(client *ethclient.Client, hash common.Hash) (uint, error) {
	count, err := client.TransactionCount(context.Background(), hash)
	if err != nil {
		log.Warnf("check account err: %v", err)
		return 0, err
	}
	return count, nil
}
