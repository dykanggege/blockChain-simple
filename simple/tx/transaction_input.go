package tx

import (
	"blockChain/simple/wallet"
	"bytes"
)

type TXInput struct {
	//该输入原来存在于哪个交易，id
	Txid []byte
	//交易的第几个输出，由交易 id 和 输出序号 共同锁定
	Vout int
	//用来解锁输出的签名
	Signature []byte
	//输入用户的公钥
	PubKey []byte
}

func (t *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(t.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
