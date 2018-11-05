package tx

import (
	"blockChain/simple/util"
	"bytes"
)

type TXOutput struct {
	//币值
	Value int
	//接受用户的公钥　hash 值
	PubKeyHash []byte
}

//将输出锁定在某个账户上
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := util.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

//该输出是否和公钥锁定
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}
