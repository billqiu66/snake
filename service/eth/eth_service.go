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

	//// 监听新区块的产生：url - ws://127.0.0.1:7545
	//go func() {
	//	headers := make(chan *types.Header)
	//	sub, err := client.SubscribeNewHead(context.Background(), headers)
	//	if err != nil {
	//		log.Warnf("subscribe blockchain err: %v", err)
	//		panic(err)
	//	}
	//	for {
	//		select {
	//		case err := <-sub.Err():
	//			log.Warnf("subscribe blockchain err: %v", err)
	//		case header := <-headers:
	//			log.Infof("[New Block]: %v", header.Hash().Hex())
	//		}
	//	}
	//
	//}()

}
