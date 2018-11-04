package tx

import (
	"blockChain/simple/util"
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

//币基交易，用于产生挖矿奖励
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			panic(err)
		}
		data = fmt.Sprintf("%x", randData)
	}

	txin := TXInput{[]byte{}, -1, nil, []byte(data)}

}

// 确定它是币基交易，用于产生挖矿奖励
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

//序列化交易
func (tx Transaction) Serialize() []byte {
	content := bytes.Buffer{}
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(tx)
	util.ErrLogPanic(err)
	return content.Bytes()
}

//反序列化交易
func DeserializeTranscation(data []byte) Transaction {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	tx := Transaction{}
	err := decoder.Decode(tx)
	util.ErrLogPanic(err)
	return tx
}

//求交易的　hash 值
func (tx *Transaction) Hash() []byte {
	txCopy := *tx
	txCopy.ID = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}
