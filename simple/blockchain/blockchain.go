package blockchain

import (
	"blockChain/simple/block"
	"blockChain/simple/tx"
	"blockChain/simple/util"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"log"
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

//创建一个从未存在的区块链
func CreateBlockchain(address, node string) *BlockChain {
	//确定区块链存储文件是否已存在
	if _, err := os.Stat(dbFile + node); !os.IsNotExist(err) {
		fmt.Println("区块链存储文件已存在")
		os.Exit(1)
	}
	bc := new(BlockChain)

	basetx := tx.NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := block.NewGenersisBlock(basetx)

	db, err := bolt.Open(dbFile+node, 0600, nil)
	util.ErrLogPanic(err)
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blockBucket))
		util.ErrLogPanic(err)
		err = b.Put(genesis.Hash, genesis.Serialize())
		util.ErrLogPanic(err)
		err = b.Put([]byte("l"), genesis.Hash)
		util.ErrLogPanic(err)
		bc.LastHash = genesis.Hash
		return nil
	})
	util.ErrLogPanic(err)
	bc.DB = db
	return bc
}

// 从一个已有的文件里，初始化区块链
func NewBlockchain(nodeID string) *BlockChain {
	if _, err := os.Stat(dbFile + nodeID); !os.IsExist(err) {
		util.ErrLogPanic(err)
	}
	db, err := bolt.Open(dbFile+nodeID, 0600, nil)
	util.ErrLogPanic(err)

	bc := new(BlockChain)

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lasthash := b.Get([]byte("l"))
		bc.LastHash = lasthash
		return nil
	})
	bc.DB = db
	return bc
}

// 传入区块数据，向区块链中添加一个区块
func (bc *BlockChain) AddBlock(bk *block.Block) {
	err := bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))

		blockIn := b.Get([]byte(bk.Hash))
		if blockIn != nil {
			return nil
		}
		blockdata := bk.Serialize()
		err := b.Put([]byte(bk.Hash), blockdata)
		util.ErrLogPanic(err)

		lasthash := b.Get([]byte("l"))
		lastblockdata := b.Get([]byte(lasthash))
		lastblock := block.DeserializeBlock(lastblockdata)

		if bk.Height > lastblock.Height {
			err := b.Put([]byte("l"), bk.Hash)
			util.ErrLogPanic(err)
			bc.LastHash = bk.Hash
		}
		return nil
	})
	util.ErrLogPanic(err)
}

func (bc *BlockChain) Iterator() *BlockChainInterator {
	return &BlockChainInterator{CurrentHash: bc.LastHash, DB: bc.DB}
}

func (bc *BlockChain) FindTranscation(ID []byte) (*tx.Transaction, error) {
	iterator := bc.Iterator()
	for iterator.Has() {
		bk := iterator.Next()
		for txix, _ := range bk.Transcations {
			if bytes.Compare(ID, bk.Transcations[txix].ID) == 0 {
				return bk.Transcations[txix], nil
			}
		}
	}
	return &tx.Transaction{}, errors.New("没有找到区块")
}

func (bc *BlockChain) FindUTXO() map[string]tx.TXOutputs {

}

//得到最近一次区块的高度
func (bc *BlockChain) GetBestHeight() int {
	var lastBlock *block.Block

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = block.DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

func (bc *BlockChain) GetBlock(blockhash []byte) (*block.Block, error) {
	iterator := bc.Iterator()
	for iterator.Has() {
		bk := iterator.Next()
		if bytes.Compare(bk.Hash, blockhash) == 0 {
			return bk, nil
		}
	}
	return nil, errors.New("找不到对应的区块")
}

func (bc *BlockChain) GetBlockHashs() [][]byte {
	blockHashs := make([][]byte, 0)
	iterator := bc.Iterator()
	for iterator.Has() {
		bk := iterator.Next()
		blockHashs = append(blockHashs, bk.Hash)
	}
	return blockHashs
}

func (bc *BlockChain) SignTranscation(t *tx.Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := map[string]tx.Transaction{}

	for _, vin := range t.Vin {
		prit, err := bc.FindTranscation(vin.Txid)
		util.ErrLogPanic(err)
		prevTXs[hex.EncodeToString(prit.ID)] = *prit
	}

	t.Sign(privKey, prevTXs)
}

func (bc *BlockChain) VerifyTranscation(t *tx.Transaction) bool {
	prevTXs := map[string]tx.Transaction{}

	for _, vin := range t.Vin {
		prevt, err := bc.FindTranscation(vin.Txid)
		util.ErrLogPanic(err)
		prevTXs[hex.EncodeToString(prevt.ID)] = *prevt
	}

	return t.Verify(prevTXs)
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
	bci.CurrentHash = bl.PrevBlockHash
	return bl
}
