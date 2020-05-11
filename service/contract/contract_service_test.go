package contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"testing"
)

func TestDeployContract(t *testing.T) {

	client, _ := ethclient.Dial("http://127.0.0.1:7545")
	pri := "c2383178616a9dfddd51d56fba9a7e98f05097d710414670532f10d56e7f744c"
	name := "gbcom"
	symbol := "gc"
	suply := big.NewInt(10000000000)
	decimal := uint8(18)

	address, tx, _ := DeployContract(pri, name, symbol, client, suply, decimal)
	println("contract addrss: ", address)
	println("tx: ", tx)
}

func TestLoadContract(t *testing.T) {

	address := "0x3F57bDDEF43eF4abDfC0dfE1852249cb62bedeD4"
	client, _ := ethclient.Dial("http://127.0.0.1:7545")

	Contracter, err := LoadContract(address, client)
	if err != nil {
		println(err)
	}
	name, err := Contracter.Name(&bind.CallOpts{})
	if err != nil {
		println(err)
	}
	println("name: ", name)
}
