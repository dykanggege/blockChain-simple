package tx

type TXOutput struct {
	//币值
	Value int
	//接受用户的公钥　hash 值
	PubKeyHash []byte
}
