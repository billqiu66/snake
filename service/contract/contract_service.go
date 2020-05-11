package contract

import (
	"github.com/1024casts/snake/service/eth"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
)

var Contracter *Supercoin
var SmartContractAddressFile = "SmartContractAddress"

func Init() {
	if fileexist(SmartContractAddressFile) {
		content, _ := ioutil.ReadFile(SmartContractAddressFile)
		address := string(content)
		log.Infof("smart contract has deployed, address: %s", address)
		var err error
		Contracter, err = LoadContract(address, eth.Client)
		if err != nil {
			panic(err)
		}
		log.Infof("load smart contract success.")
	} else {
		pri := viper.GetString("eth.privatekey")
		coin := viper.GetStringMapString("eth.coin")
		temp_64, _ := strconv.ParseInt(coin["initialsupply"], 10, 64)
		initialsupply := big.NewInt(temp_64)
		decimals, _ := strconv.ParseUint(coin["decimals"], 10, 8)

		address, _, err := DeployContract(pri, coin["name"], coin["symbol"], eth.Client, initialsupply, uint8(decimals))
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(SmartContractAddressFile, []byte(address), 0664)
	}
}

func LoadContract(address string, client *ethclient.Client) (*Supercoin, error) {
	addr := common.HexToAddress(address)
	instance, err := NewSupercoin(addr, client)
	if err != nil {
		log.Warnf("load smart contract fail: %v", err)
		return nil, err
	}
	return instance, nil
}

func DeployContract(pri, name, symbol string, client *ethclient.Client, suply *big.Int, decimals uint8) (string, string, error) {
	privateKey, err := crypto.HexToECDSA(pri)
	if err != nil {
		log.Warnf("deploy smart contract err: %v", err)
		return "", "", err
	}

	//publicKey := privateKey.Public()
	//publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	//if !ok {
	//	log.Warn("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	//	return "", "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	//}
	//
	//fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	//nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	//if err != nil {
	//	log.Warnf("deploy smart contract err: %v", err)
	//	return "", "", nil
	//}
	//
	//gasPrice, err := client.SuggestGasPrice(context.Background())
	//if err != nil {
	//	log.Warnf("deploy smart contract err: %v", err)
	//	return "", "", nil
	//}

	auth := bind.NewKeyedTransactor(privateKey)
	//auth.Nonce = big.NewInt(int64(nonce))
	//auth.Value = big.NewInt(0)     // in wei
	//auth.GasLimit = uint64(300000) // in units
	//auth.GasPrice = gasPrice

	address, tx, instance, err := DeploySupercoin(auth, client, suply, name, symbol, decimals)
	if err != nil {
		log.Warnf("deploy smart contract err: %v", err)
		return "", "", nil
	}
	tokenName, err := instance.Name(&bind.CallOpts{Pending: true})
	if err != nil {
		println(err)
	}
	println("name: ", tokenName)
	Contracter = instance
	//log.Warnf("contract address: %v", address.Hex())
	//log.Warnf("tx: %v", tx.Hash().Hex())
	return address.Hex(), tx.Hash().Hex(), nil
}

func fileexist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
