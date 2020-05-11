package contract

import (
	"github.com/1024casts/snake/service/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"testing"
)

func TestBalanceAt(t *testing.T) {
	//address := "0x6B0f066eC74f89d32f5A9937CFB9457E64359Ae7"
	address := "0xFFeF3BEAFaDD2067aAFf431F0cF9eA7e25546C15"
	contractAddress := "0x93bb81EF84B992D9C02Cdc9eCd457C0B64a5ECDC"
	client, _ := ethclient.Dial("http://127.0.0.1:7545")
	instance, _ := contract.NewSupercoin(common.HexToAddress(contractAddress), client)

	balance, _ := BalanceAt(instance, address)
	log.Printf("balance: %v", balance)
}

func TestTransferAt(t *testing.T) {
	private := "c2383178616a9dfddd51d56fba9a7e98f05097d710414670532f10d56e7f744c"
	address := "0xFFeF3BEAFaDD2067aAFf431F0cF9eA7e25546C15"
	contractAddress := "0x93bb81EF84B992D9C02Cdc9eCd457C0B64a5ECDC"
	client, _ := ethclient.Dial("http://127.0.0.1:7545")

	tx, err := TransferAt(client, private, address, contractAddress)
	log.Printf("tx:%v, err:%v", tx, err)
}
