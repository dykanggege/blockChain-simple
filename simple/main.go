package main

import (
	"blockChain/simple/cli"
	"fmt"
)

func main() {
	c := cli.Cli{}
	for {
		var cmd string
		fmt.Scanln(&cmd)
		c.Run(cmd)
	}
}
