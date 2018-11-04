//真　区块
package block

import (
	"blockChain/simple/merkleTree"
	"blockChain/simple/pow"
	"blockChain/simple/tx"
	"blockChain/simple/util"
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

// 一个区块的结构
type Block struct {
	// 创建区块的时间戳
	TimeStamp int64
	// 区块中的包含的交易
	Transcations []*tx.Transaction
	// 前一个区块的 hash 值
	PrevBlockHash []byte
	//幸运数字
	LuckNum int64
	// 以上面的字段内容，计算得到的 hash 值
	Hash []byte
	//区块高度
	Height int
}

// 创建一个新区块
func NewBlock(trans []*tx.Transaction, prevHash []byte, height int) *Block {
	b := &Block{TimeStamp: time.Now().Unix(), Transcations: trans, PrevBlockHash: prevHash, Hash: []byte{}, LuckNum: 0, Height: height}
	// 挖矿，得到幸运数字，和区块的 hash 值
	power := pow.NewProofOfWork(b)
	luck, hash := power.Run()

	b.LuckNum = luck
	b.Hash = hash
	return b
}

//创建　创世区块
func NewGenersisBlock(coinbase *tx.Transaction) *Block {
	return NewBlock([]*tx.Transaction{}, []byte{}, 0)
}

//　得到所有交易的　merkleTree　的根的 hash 值
func (b *Block) HashTranscations() []byte {
	var txs [][]byte

	for i, _ := range b.Transcations {
		txs = append(txs, b.Transcations[i].Serialize())
	}

	mTree := merkleTree.NewMerkleTree(txs)

	return mTree.RootNode.Data
}

// 使用 go 内置的 gob 编码方式，将其编码为 byte 类型，用于存储
func (b *Block) Serialize() []byte {
	var data bytes.Buffer
	encoder := gob.NewEncoder(&data)
	err := encoder.Encode(b)
	util.ErrLogPanic(err)
	return data.Bytes()
}

// 将数据解码，得到区块
func DeserializeBlock(data []byte) *Block {
	b := new(Block)
	decoder := gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(b)
	return b
}

func (b *Block) Print() {
	fmt.Printf("TimeStamp: %d\n", b.TimeStamp)
	fmt.Printf("Data: %s\n", b.Transcations)
	fmt.Printf("PrevBlockHash: %x\n", b.PrevBlockHash)
	fmt.Printf("LuckNum: %d\n", b.LuckNum)
	fmt.Printf("Hash: %x\n", b.Hash)
	power := pow.NewProofOfWork(b)
	fmt.Printf("Validate: %v\n", power.Validate())
	fmt.Println()
}
