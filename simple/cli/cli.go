package cli

import (
	"blockChain/simple/util"
	"flag"
	"fmt"
	"os"
)

type Cli struct {
	flags map[string]*flag.FlagSet
}

func New() *Cli {
	c := Cli{}
	c.flags = make(map[string]*flag.FlagSet)
	return &c
}

func (c *Cli) Run() {
	//验证传入参数
	c.validateArgs()
	//注册要输入的命令
	f := c.registerCmd("createwallet")

	if f, ok := c.flags[os.Args[1]]; ok {
		err := f.Parse(os.Args[2:])
		util.ErrLogPanic(err)
	} else {
		c.printUsage()
		os.Exit(1)
	}
}

func (c Cli) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  reindexutxo - Rebuilds the UTXO set")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT -mine - Send AMOUNT of coins from FROM address to TO. Mine on the same node, when -mine is set.")
	fmt.Println("  startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining")
}

func (c *Cli) validateArgs() {
	if len(os.Args) < 2 {
		c.printUsage()
		os.Exit(1)
	}
}

//将 cmd 注册在 cmds 中
func (c *Cli) registerCmd(cmd string) *flag.FlagSet {
	f := flag.NewFlagSet(cmd, flag.ExitOnError)
	c.flags[cmd] = f
	return f
}
