package main

import (
	"github.com/jeonjonghyeok/coin/blockchain"
	"github.com/jeonjonghyeok/coin/cli"
	"github.com/jeonjonghyeok/coin/db"
)

func main() {
	blockchain.Blockchain()
	defer db.Close()
	cli.Start()
}
