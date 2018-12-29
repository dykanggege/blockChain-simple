package net

import (
	"blockChain/simple/block"
	"blockChain/simple/tx"
	"blockChain/simple/util"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
)

//使用的网络协议
const protocol = "tcp"

//网络服务的版本号
const nodeVersion = 1

//命令字节长度限制
const commandLength = 12

//我的节点地址
var nodeAddress string

//挖矿节点地址
var miningAddress string

//已知节点
var knownNodes = []string{"localhost:3000"}
var blocksInTransit = [][]byte{}
var mempool = make(map[string]tx.Transaction)

type Addr struct {
	AddrList []string
}

type Block struct {
	AddrFrom string
	Block    []byte
}

//block 数据来源的节点地址
type Getblocks struct {
	AddrFrom string
}

type Getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type Inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

type Tx struct {
	AddFrom     string
	Transcation []byte
}

type Version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

//向所有已知节点请求区块
func RequestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
}

//向一个地址发送所有我知道的地址
func SendAddr(address string) {
	nodes := Addr{knownNodes}
	nodes.AddrList = append(nodes.AddrList, nodeAddress)
	payload := gobEncode(nodes)
	request := append(commandToBytes("addr"), payload...)

	sendData(address, request)
}

//向某个地址发送区块
func SendBlock(addr string, b *block.Block) {
	data := Block{nodeAddress, b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

//发送所有的货存
func SendInv(addr, kind string, items [][]byte) {
	inventory := Inv{nodeAddress, kind, items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(addr, request)
}

func SendGetBlocks(addr string) {
	payload := gobEncode(Getblocks{nodeAddress})

}

//将指令变为字节便于传输，固定指令的字节长度
func commandToBytes(command string) []byte {
	bs := make([]byte, commandLength)
	for i, c := range command {
		bs[i] = byte(c)
	}
	return bs[:]
}

//清除无效字节，翻译成指令字符串
func bytesToCommand(bs []byte) string {
	var command []byte

	for _, b := range bs {
		if b != 0x0 {
			command = append(command, b)
		}
	}
	return string(command)
}

func ExtractCommand(request []byte) []byte {
	return request[:commandLength]
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

func sendGetBlocks(address string) {
	//将命令和数据封装在一起发送
	payload := gobEncode(Getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Println("%s 不可用", addr)
		//如果发送的节点不可用，清除该节点
		for i, n := range knownNodes {
			if n == addr {
				knownNodes = append(knownNodes[:i], knownNodes[i+1:]...)
				return
			}
		}
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	util.ErrLogPanic(err)
}
