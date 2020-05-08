package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

var Client *ethclient.Client

func Init() {
	ethUrl := viper.GetString("eth.url")

	client, err := ethclient.Dial(ethUrl)
	if err != nil {
		log.Warnf("Dial to eth err: %v", err)
		panic(err)
	}
	Client = client
}
