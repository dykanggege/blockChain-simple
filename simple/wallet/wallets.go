// 钱包　wallets，可以包含多个账户 wallet，封装了钱包常用的操作
package wallet

import (
	"blockChain/simple/util"
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

//钱包文件存放的地方
const walletFile = ""

type Wallets struct {
	Wallets map[string]*Wallet
}

//以钱包目录为根目录，从某个文件里加载一个钱包
func NewWallets(nodeID string) (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.loadFromFile(nodeID)

	return &wallets, err
}

//创建一个账号，返回地址
func (w *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := string(wallet.GetAddress())

	w.Wallets[address] = wallet

	return address
}

//得到钱包里所有的地址
func (w *Wallets) GetAddresses() []string {
	s := make([]string, 0)
	for key, _ := range w.Wallets {
		s = append(s, key)
	}
	return s
}

func (w *Wallets) GetWallet(addr string) *Wallet {
	return w.Wallets[addr]
}

//以钱包文件夹为基，将钱包存入到某个文件中
func (w Wallets) SaveToFile(nodeID string) {
	var conetnt bytes.Buffer
	encoder := gob.NewEncoder(&conetnt)
	err := encoder.Encode(w)
	util.ErrLogPanic(err)

	err = ioutil.WriteFile(walletFile+nodeID, conetnt.Bytes(), 0644)
	util.ErrLogPanic(err)
}

//从一个文件中加载钱包
func (w *Wallets) loadFromFile(nodeID string) error {
	//从文件夹的第　node 个节点加载钱包文件
	walletFile := fmt.Sprintf(walletFile, nodeID)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	util.ErrLogPanic(err)

	var wallets Wallets
	gob.Register(elliptic.P256())
	decode := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decode.Decode(&wallets)
	util.ErrLogPanic(err)

	w.Wallets = wallets.Wallets

	return nil
}
