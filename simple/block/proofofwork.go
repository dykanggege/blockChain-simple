package block

import (
	"blockChain/simple/util"
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"runtime"
	"sync"
)

// 工作难度
const TARGETBITS = 20

// 要生成的 hash 值长度
const HASHLENGTH = 256

// 幸运数字的大小限制
const MAXLUCKNUM = math.MaxInt64

// 工作量证明
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target = target.Lsh(target, uint(HASHLENGTH-TARGETBITS))

	return &ProofOfWork{block: b, target: target}
}

// 得到区块的 bytes 数据,不包含 luckNum
func (pow *ProofOfWork) prepareData() []byte {
	return bytes.Join(
		[][]byte{
			pow.block.Data,
			pow.block.PrevHash,
			util.IntToHex(pow.block.TimeStamp),
		}, []byte{})
}

// 挖矿
func (pow *ProofOfWork) Run() (int64, []byte) {
	var hash [sha256.Size]byte
	var luckNum int64
	data := pow.prepareData()

	lock := sync.Mutex{}
	cpuNum := int64(runtime.NumCPU())
	stop := make(chan bool)
	clo := false
	var interval int64 = 10000
	var i int64
	for i = 0; i < cpuNum; i++ {
		go func() {
			var hashint big.Int
			var luck int64
			for start := i * interval; start < MAXLUCKNUM; start += cpuNum * interval {
				for luck = start; luck < start+interval; luck++ {
					lock.Lock()
					if clo {
						return
					}
					lock.Unlock()

					newdata := append(data, util.IntToHex(luck)...)
					h := sha256.Sum256(newdata)
					hashint.SetBytes(h[:])

					if hashint.Cmp(pow.target) == -1 {
						lock.Lock()
						defer lock.Unlock()
						hash = h
						luckNum = luck
						clo = true
						stop <- true
					}
				}
			}
		}()
	}

	<-stop
	return luckNum, hash[:]
}

// 检验幸运数字是否正确
func (pow *ProofOfWork) Validate() bool {
	var hashint big.Int
	data := append(pow.prepareData(), util.IntToHex(pow.block.LuckNum)...)
	hash := sha256.Sum256(data)
	hashint.SetBytes(hash[:])
	return hashint.Cmp(pow.target) == -1
}
