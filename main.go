package main

import (
	"github.com/jeonjonghyeok/coin/blockchain"
	"github.com/jeonjonghyeok/coin/cli"
)

func main() {
	blockchain.Blockchain()
	cli.Start()
}
