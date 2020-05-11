package account

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/1024casts/snake/handler"
	"github.com/1024casts/snake/pkg/errno"
	"github.com/1024casts/snake/service/eth"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func CheckAddress(c *gin.Context) {
	address := c.PostForm("address")
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	if re.MatchString(address) {
		address := common.HexToAddress(address)
		bytecode, err := eth.Client.CodeAt(context.Background(), address, nil) // nil is latest block
		if err != nil {
			log.Warnf("check address err: %v", err)
			handler.SendResponse(c, errno.ErrAccountCheck, nil)
			return
		}

		isContract := len(bytecode) > 0
		if isContract {
			handler.SendResponse(c, nil, "the address is the contract address")
		} else {
			handler.SendResponse(c, nil, "the address is the common address")
		}

	} else {
		log.Warnf("invalid eth address: %s", address)
		handler.SendResponse(c, errno.ErrAccountInvalid, nil)
	}
}

func NewAccount(c *gin.Context) {

	// 1. 生成随机私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Warnf("crypto GenerateKey err: %v", err)
		handler.SendResponse(c, errno.ErrAccount, nil)
		return
	}
	// 2. 将私钥转换成字节
	privateKeyBytes := crypto.FromECDSA(privateKey)
	// 3. 十六进制编码后，得到私钥
	pri := hexutil.Encode(privateKeyBytes)[2:]
	// 4. 从私钥中获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Warn("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		handler.SendResponse(c, errno.ErrAccount, nil)
		return
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	pub := hexutil.Encode(publicKeyBytes)[4:]
	// 5. 从公钥中获取地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	res := map[string]string{
		"pri":     pri,
		"pub":     pub,
		"address": address,
	}

	handler.SendResponse(c, nil, res)
	return

}

func NewKeyStore(c *gin.Context) {
	password := c.PostForm("password")
	if password == "" {
		log.Warn("argument error, miss password!")
		handler.SendResponse(c, errno.ErrAccountPassword, nil)
		return
	}

	ks := keystore.NewKeyStore("./tmp", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Warnf("create keystore error: %v", err)
		return
	}

	address := account.Address.Hex()
	log.Info(address) // 0x20F8D42FB0F667F2E53930fed426f225752453b3
	//var content string
	var keystore map[string]interface{}
	err = filepath.Walk("./tmp",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Warn(err.Error())
				return err
			}

			if info.IsDir() {
				return nil
			}
			log.Infof("file name: %s, %s", info.Name(), strings.ToLower(address[2:]))
			ok := strings.HasSuffix(info.Name(), strings.ToLower(address[2:]))
			if ok {
				file, err := ioutil.ReadFile("./tmp/" + info.Name())
				if err != nil {
					log.Warnf("open file err: %v", err)
				}
				//content = string(file)
				if err := json.Unmarshal(file, &keystore); err == nil {
					log.Infof("%v", keystore)
				} else {
					log.Warnf("%v", err)
				}
				return nil
			}
			return nil

		})

	if err != nil {
		log.Warnf("create keystore err: %v", err)
		handler.SendResponse(c, errno.ErrAccountPassword, nil)
		return
	}
	res := map[string]interface{}{"keystore": keystore, "address": address}
	handler.SendResponse(c, nil, res)
}

func Verify(c *gin.Context) {

}
