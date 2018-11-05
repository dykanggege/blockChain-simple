package blockchain

import (
	"blockChain/simple/block"
	"blockChain/simple/util"
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

const dbFile = "data.DB"
const blockBucket = "blocks"
const genesisCoinbaseData = "这是最好的时代，也是最坏的时代"

type BlockChain struct {
	// 存储好多好多区块，并连接成一串
	// 通过前后 hash 值关联自发连接
	DB       *bolt.DB
	LastHash []byte
}

func CreateBlockchain(address, node string) *BlockChain {
	//确定区块链存储文件是否已存在
	if _, err := os.Stat(dbFile + node); !os.IsNotExist(err) {
		fmt.Println("区块链存储文件已存在")
		os.Exit(1)
	}

}

// 创建初始区块
func newGenesisBlock() *block.Block {
	return block.NewBlock([]byte("2018/10/19/11:16 我坐在火车上，百无聊赖的写下这行代码"), []byte{})
}

// 初始化区块链
func New() *BlockChain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	util.ErrLogPanic(err)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))

		if b == nil {
			bl := newGenesisBlock()
			b, err := tx.CreateBucket([]byte(blockBucket))
			util.ErrLogPanic(err)
			err = b.Put(bl.Hash, bl.Serialize())
			util.ErrLogPanic(err)
			err = b.Put([]byte("l"), bl.Hash)
			util.ErrLogPanic(err)
			tip = bl.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	return &BlockChain{LastHash: tip, DB: db}
}

// 传入区块数据，向区块链中添加一个区块
func (bc *BlockChain) AddBlock(data string) {
	var lasthash []byte
	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lasthash = b.Get([]byte("l"))
		return nil
	})
	util.ErrLogPanic(err)
	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		bl := block.NewBlock([]byte(data), lasthash)
		b.Put(bl.Hash, bl.Serialize())
		b.Put([]byte("l"), bl.Hash)
		bc.LastHash = bl.Hash
		return nil
	})
}

func (bc *BlockChain) Iterator() *BlockChainInterator {
	return &BlockChainInterator{CurrentHash: bc.LastHash, DB: bc.DB}
}

type BlockChainInterator struct {
	CurrentHash []byte
	DB          *bolt.DB
}

func (bci BlockChainInterator) Has() bool {
	return len(bci.CurrentHash) != 0
}

func (bci *BlockChainInterator) Next() *block.Block {
	var bl *block.Block
	bci.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		encodeblock := b.Get(bci.CurrentHash)
		bl = block.DeserializeBlock(encodeblock)
		return nil
	})
	bci.CurrentHash = bl.PrevHash
	return bl
}
