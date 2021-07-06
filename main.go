package main

import (
	"github.com/jeonjonghyeok/coin/blockchain"
	"github.com/jeonjonghyeok/coin/cli"
	"github.com/jeonjonghyeok/coin/wallet"
)

func main() {
	wallet.Wallet()
	blockchain.Blockchain()
	cli.Start()
}
