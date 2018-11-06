package tx

import (
	"blockChain/simple/block"
	"blockChain/simple/blockchain"
	"blockChain/simple/util"
	"encoding/hex"
	"github.com/boltdb/bolt"
)

//存储所有未被使用的输出
const utxoBucket = "chainstate"

//存储所有未被使用的输出
//存储形式 key:TXID val:交易中所有的输出
type UTXOset struct {
	Blockchain *blockchain.BlockChain
}

//寻找用户可用于支付的输出
//将多个小的余额拼凑为一个大的余额
func (u UTXOset) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, map[string][]int) {
	unspendOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.DB

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeTXOutPutes(v)

			for outIndex, out := range outs.Outputs {
				if out.IsLockedWithKey(pubkeyHash) && accumulated < amount {
					accumulated += out.Value
					unspendOutputs[txID] = append(unspendOutputs[txID], outIndex)
				}
			}

		}
		return nil
	})
	util.ErrLogPanic(err)

	return accumulated, unspendOutputs
}

//查找余额
func (u UTXOset) FindUTXO(pubkeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.DB

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeTXOutPutes(v)
			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubkeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}
		return nil
	})
	util.ErrLogPanic(err)

	return UTXOs
}

//返回所有仍有可用输出的交易数量
func (u UTXOset) CountTransactions() int {
	db := u.Blockchain.DB
	count := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			count++
		}
		return nil
	})
	util.ErrLogPanic(err)
	return count
}

//更新可用于支付的输出表
func (u UTXOset) Update(block *block.Block) {
	db := u.Blockchain.DB

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transcations {
			if !tx.IsCoinbase() {
				//遍历区块中所有的输入
				for _, vin := range tx.Vin {
					updateOuts := TXOutputs{}
					outsBytes := b.Get(vin.Txid)
					outs := DeserializeTXOutPutes(outsBytes)

					for outIndex, out := range outs.Outputs {
						if outIndex != vin.Vout {
							updateOuts.Outputs = append(updateOuts.Outputs, out)
						}
					}

					if len(updateOuts.Outputs) == 0 {
						err := b.Delete(vin.Txid)
						util.ErrLogPanic(err)
					} else {
						err := b.Put(vin.Txid, updateOuts.Serialize())
						util.ErrLogPanic(err)
					}
				}
			}

			newOutputs := TXOutputs{}
			for _, out := range tx.Vout {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}
			err := b.Put(tx.ID, newOutputs.Serialize())
			util.ErrLogPanic(err)
		}
		return nil
	})
	util.ErrLogPanic(err)
}
