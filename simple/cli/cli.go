package cli

import (
	"fmt"
)

type Cli struct{}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data 添加区块的数据")
}

func (c *Cli) Run(cmd string) {
	option := ""
	fmt.Sscan(cmd, &option)

	switch option {
	case "addblock":
	default:
		fmt.Println("无效的命令!!!")
		printUsage()
		return
	}
}
