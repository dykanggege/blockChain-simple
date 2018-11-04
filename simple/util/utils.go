package util

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

// 将 int64 类型数据转为 []byte
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	ErrLogPanic(err)
	return buff.Bytes()
}

// 翻转一个　[]byte
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func ErrLogPanic(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func FileExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
