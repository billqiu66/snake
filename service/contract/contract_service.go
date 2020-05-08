package contract

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/1024casts/snake/handler/smartcontract"
	"github.com/1024casts/snake/service/eth"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
)

var Contracter *smartcontract.Supercoin
var SmartContractAddressFile = "SmartContractAddress"

func Init() {
	if fileexist(SmartContractAddressFile) {
		content, _ := ioutil.ReadFile(SmartContractAddressFile)
		address := string(content)
		log.Infof("smart contract has deployed, address: %s", address)
		err := loadContract(address)
		if err != nil {
			panic(err)
		}
		log.Infof("load smart contract success.")
	} else {
		address, _, err := deployContract()
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(SmartContractAddressFile, []byte(address), 0664)
	}
}

func loadContract(address string) error {
	addr := common.HexToAddress(address)
	instance, err := smartcontract.NewSupercoin(addr, eth.Client)
	if err != nil {
		log.Warn("load smart contract success.")
		return err
	}
	Contracter = instance
	return nil
}

func deployContract() (string, string, error) {
	pri := viper.GetString("eth.privatekey")
	privateKey, err := crypto.HexToECDSA(pri)
	if err != nil {
		log.Warnf("deploy smart contract err: %v", err)
		return "", "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Warn("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return "", "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := eth.Client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Warnf("deploy smart contract err: %v", err)
		return "", "", nil
	}

	gasPrice, err := eth.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Warnf("deploy smart contract err: %v", err)
		return "", "", nil
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	coin := viper.GetStringMapString("eth.coin")
	temp_64, _ := strconv.ParseInt(coin["initialsupply"], 10, 64)
	initialsupply := big.NewInt(temp_64)
	decimals, _ := strconv.ParseUint(coin["decimals"], 10, 8)
	address, tx, instance, err := smartcontract.DeploySupercoin(auth, eth.Client, initialsupply, coin["name"], coin["symbol"], uint8(decimals))
	if err != nil {
		log.Warnf("deploy smart contract err: %v", err)
		return "", "", nil
	}
	Contracter = instance
	log.Warnf("contract address: %v", address.Hex())
	log.Warnf("tx: %v", tx.Hash().Hex())
	return address.Hex(), tx.Hash().Hex(), nil
}

func fileexist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
