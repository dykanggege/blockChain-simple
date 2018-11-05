package tx

import (
	"blockChain/simple/util"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
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

//求交易的　sha256 hash 值
func (tx *Transaction) Hash() []byte {
	txCopy := *tx
	txCopy.ID = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

//对交易做签名
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	//币基交易不需要签名，大家承认即可
	if tx.IsCoinbase() {
		return
	}

	for i, _ := range tx.Vin {
		if _, ok := prevTXs[hex.EncodeToString(tx.Vin[i].Txid)]; !ok {
			log.Panic("错误: 输入的前一个交易并不存在，该输入可能是伪造的")
		}
	}

	txCopy := tx.TrimmedCopy()

	for index, val := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(val.Txid)]
		txCopy.Vin[index].Signature = nil
		txCopy.Vin[index].PubKey = prevTx.Vout[val.Vout].PubKeyHash

		dataToSign := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
		util.ErrLogPanic(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[index].Signature = signature
		txCopy.Vin[index].PubKey = nil
	}
}

//　验证签名是否有效
func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for i, _ := range tx.Vin {
		if _, ok := prevTxs[hex.EncodeToString(tx.Vin[i].Txid)]; !ok {
			log.Panic("输入的源交易不存在")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for index, val := range tx.Vin {
		prevTX := prevTxs[hex.EncodeToString(val.Txid)]
		txCopy.Vin[index].Signature = nil
		txCopy.Vin[index].PubKey = prevTX.Vout[val.Vout].PubKeyHash

		r := big.Int{}
		s := big.Int{}
		siglen := len(val.Signature)
		r.SetBytes(val.Signature[:(siglen / 2)])
		s.SetBytes(val.Signature[(siglen / 2):])

		x := big.Int{}
		y := big.Int{}
		keylen := len(val.PubKey)
		x.SetBytes(val.PubKey[:(keylen / 2)])
		y.SetBytes(val.PubKey[(keylen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Vin[index].PubKey = nil
	}
	return true
}

// 确定它是币基交易，用于产生挖矿奖励
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))

	for i, input := range tx.Vin {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.Txid))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Vout {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Value))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}

//对交易做分离的拷贝
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

//创建一个币基交易
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		util.ErrLogPanic(err)
		data = fmt.Sprintf("%x", randData)
	}

	txin := TXInput{Txid: []byte{}, Vout: -1, Signature: nil, PubKey: []byte(data)}

}
