package contract

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/service/contract"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"math/big"
)

func Check(c *gin.Context) {
	name, err := contract.Contracter.Name(nil)

	if err != nil {
		log.Warnf("get contract info fail: %v", err)
	}

	symbol, _ := contract.Contracter.Symbol(nil)
	initialsupply, _ := contract.Contracter.TotalSupply(nil)
	decimals, _ := contract.Contracter.Decimals(nil)

	res := map[string]interface{}{
		"name":          name,
		"symbol":        symbol,
		"initialsupply": initialsupply,
		"decimals":      decimals,
	}
	handler.SendResponse(c, nil, res)
}

func Balance(c *gin.Context) {
	address := c.PostForm("address")
	if address == "" {
		log.Warnf("miss argument(s)")
		handler.SendResponse(c, errno.ErrContractMissArg, nil)
		return
	}
	addr := common.HexToAddress(address)

	balance, err := contract.Contracter.BalanceOf(nil, addr)
	if err != nil {
		log.Warnf("get balance err: %v", err)
		handler.SendResponse(c, errno.ErrContractBalance, nil)
		return
	}
	req := map[string]interface{}{
		"address": address,
		"balance": balance,
	}
	handler.SendResponse(c, nil, req)
}

func Transfer(c *gin.Context) {
	Address := c.PostForm("address")

	if Address == "" {
		log.Warnf("miss argument(s)")
		handler.SendResponse(c, errno.ErrContractMissArg, nil)
		return
	}
	Addr := common.HexToAddress(Address)
	num := big.NewInt(10)
	tx, err := contract.Contracter.Transfer(nil, Addr, num)
	if err != nil {
		log.Warnf("transfer err: %v", err)
		handler.SendResponse(c, errno.ErrContractTransfer, nil)
		return
	}
	log.Infof("transfer tx:%v", tx.Hash().Hex())
}

func TransferAt(client *ethclient.Client, private, toAddr, tokenAddr string) (string, error) {
	privateKey, err := crypto.HexToECDSA(private)
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Warn("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return "", errors.New("transfer error")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}
	// token转账时，此处eth值要设置为0
	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}
	toAddress := common.HexToAddress(toAddr)
	tokenAddress := common.HexToAddress(tokenAddr)

	//智能合约函数名的ID
	transferFnSignature := []byte("transfer(address,uint256)")
	methodID := crypto.Keccak256(transferFnSignature)[:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	amount := new(big.Int)
	amount.SetString("100000", 10) // 10 tokens
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// 评估执行合约需要的gas
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}
	println("gaslimit: ", gasLimit)

	tx := types.NewTransaction(nonce, tokenAddress, value, uint64(43000), gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Warnf("transfer err: %v", err)
		return "", err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		//log.Warnf("transfer err: %v", err)
		println(err)
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

func BalanceAt(cont *contract.Supercoin, addr string) (*big.Int, error) {
	address := common.HexToAddress(addr)
	balance, err := cont.BalanceOf(nil, address)
	if err != nil {
		log.Warnf(" get balance err: %v", err)
		return nil, err
	}
	return balance, nil
}
