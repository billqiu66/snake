package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"testing"
)

func TestTransfer(t *testing.T) {
	client, _ := ethclient.Dial("http://127.0.0.1:7545")
	private := "c2383178616a9dfddd51d56fba9a7e98f05097d710414670532f10d56e7f744c"
	address := "0xFFeF3BEAFaDD2067aAFf431F0cF9eA7e25546C15"
	num := 1
	tx, _ := Transfer(client, private, address, int64(num))
	println("tx:", tx)
}

func TestBalance(t *testing.T) {
	client, _ := ethclient.Dial("http://127.0.0.1:7545")
	address := "0xFFeF3BEAFaDD2067aAFf431F0cF9eA7e25546C15"
	balance, _ := Balance(client, address)
	log.Printf("balance: %v", balance)
}
