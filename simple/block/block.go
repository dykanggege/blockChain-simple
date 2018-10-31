package block

import (
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
	// 区块中的数据
	Data []byte
	// 前一个区块的 hash 值
	PrevHash []byte
	//幸运数字
	LuckNum int64
	// 以上面的字段内容，计算得到的 hash 值
	Hash []byte
}


// 创建一个新区块
func NewBlock(data,prevHash []byte) *Block {
	b := &Block{TimeStamp:time.Now().Unix(),Data:data,PrevHash:prevHash}
	// 挖矿，得到幸运数字，和区块的 hash 值
	pow := NewProofOfWork(b)
	luch, hash := pow.Run()

	b.LuckNum = luch
	b.Hash = hash
	return b
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

func (b *Block)Print()  {
	fmt.Printf("TimeStamp: %d\n",b.TimeStamp)
	fmt.Printf("Data: %s\n",b.Data)
	fmt.Printf("PrevHash: %x\n",b.PrevHash)
	fmt.Printf("LuckNum: %d\n",b.LuckNum)
	fmt.Printf("Hash: %x\n",b.Hash)
	pow := NewProofOfWork(b)
	fmt.Printf("Validate: %v\n",pow.Validate())
	fmt.Println()
}

