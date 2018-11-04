package main

import (
	"blockChain/simple/cli"
	"fmt"
)

func main() {
	c := cli.Cli{}
	var cmd string
	for {
		fmt.Scanln(&cmd)
		fmt.Println(cmd)
		c.Run(cmd)
	}
}
