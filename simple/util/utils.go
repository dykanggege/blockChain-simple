package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

// 将 int64 类型数据转为 []byte
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff,binary.BigEndian,num)
	ErrLogPanic(err)
	return buff.Bytes()
}

func ErrLogPanic(err error)  {
	if err != nil{
		log.Panic(err)
	}
}

func FileExist(file string) bool {
	if _,err :=os.Stat(file); os.IsNotExist(err){
		return false
	}
	return true
}

func WaitExit()  {
	var i int
	fmt.Scan("%d",i)
	os.Exit(1)
}