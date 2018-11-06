package tx

import (
	"blockChain/simple/util"
	"bytes"
	"encoding/gob"
)

type TXOutput struct {
	//币值
	Value int
	//接受用户的公钥　hash 值
	PubKeyHash []byte
}

//将输出锁定在某个地址上
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := util.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

//检测该输出是否和公钥锁定
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

//创建交易输出，将地址变为pukhash存入
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{Value: value, PubKeyHash: nil}
	txo.Lock([]byte(address))

	return txo
}

//输出的集和
type TXOutputs struct {
	Outputs []TXOutput
}

func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(outs)
	util.ErrLogPanic(err)

	return buff.Bytes()
}

func DeserializeTXOutPutes(data []byte) TXOutputs {
	outs := new(TXOutputs)

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(outs)
	util.ErrLogPanic(err)
	return *outs
}
