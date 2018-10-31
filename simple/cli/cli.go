package cli

import (
	"blockChain/simple/blockchain"
	"fmt"
)

type Cli struct {
	
}

func (c *Cli)Run(cmd string)  {
	var option string
	fmt.Sscan(cmd,&option)

	switch option {
	case "addblock":
		addBlock(cmd)
	}
}

func addBlock(cmd string)  {
	data := ""
	fmt.Sscan(cmd,_,&data)
	bc := blockchain.New()
	bc.AddBlock(data)
}