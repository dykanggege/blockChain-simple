package wallet

import (
	"blockChain/simple/util"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

//版本号
const version = byte(0x00)

//地址校验和长度
const addressChecksumLen = 4

// 一个账户，包含一个公钥和私钥
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	pri, pub := newKeyPair()
	return &Wallet{pri, pub}
}

//生成地址（账户）
func (w Wallet) GetAddress() []byte {
	//为　pubkey　做 hash 运算
	pubkeyHash := HashPubKey(w.PublicKey)

	//插入版本号和校验和
	versiondPayload := append([]byte{version}, pubkeyHash...)
	checksum := checksum(versiondPayload)
	fullPayload := append(versiondPayload, checksum...)

	//以 base58 进行编码
	address := util.Base58Encode(fullPayload)
	return address
}

// 检查地址格式是否正确
func ValidateAddress(address string) bool {
	pubKeyHash := util.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

//　为公钥做 sha256 和　ripemd160 哈希运算
func HashPubKey(pubkey []byte) []byte {
	//先对　pubkey　做　sha256 hash 运算
	pub256 := sha256.Sum256(pubkey)

	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(pub256[:])
	if err != nil {
		log.Panic(err)
	}
	pubRipemd160 := ripemd160Hasher.Sum(nil)

	return pubRipemd160
}

//校验和
func checksum(payload []byte) []byte {
	sha1 := sha256.Sum256(payload)
	sha2 := sha256.Sum256(sha1[:])

	return sha2[:addressChecksumLen]
}

//　生成一对公钥和私钥
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	//创建一个椭圆
	curve := elliptic.P256()
	//使用椭圆和随机数生成私钥和公钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	util.ErrLogPanic(err)
	pubkey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return *privateKey, pubkey
}
