package pow

import (
	"blockChain/simple/block"
	"blockChain/simple/util"
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"runtime"
	"sync"
)

// 工作难度，不做动态调整，固定难度
const Targetbits = 20

// 要生成的 hash 值长度
const hashlength = 256

// 幸运数字的大小限制
const maxLucknum = math.MaxInt64

// 工作量证明
type ProofOfWork struct {
	block  *block.Block
	target *big.Int
}

func NewProofOfWork(b *block.Block) *ProofOfWork {
	target := big.NewInt(1)
	target = target.Lsh(target, uint(hashlength-Targetbits))

	return &ProofOfWork{block: b, target: target}
}

// 得到区块的 bytes 数据,不包含 luckNum
//包含在区块链中的区块，除了区块的数据，还有目标难度和版本号等，这里做了简化，只包含目标难度
func (pow *ProofOfWork) prepareData() []byte {
	return bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			//使用交易的 merkleTree 根哈希值作为数据填充
			pow.block.HashTranscations(),
			util.IntToHex(pow.block.TimeStamp),
			util.IntToHex(Targetbits),
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
			for start := i * interval; start < maxLucknum; start += cpuNum * interval {
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
