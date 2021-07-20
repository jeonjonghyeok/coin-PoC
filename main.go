package main

import (
	"github.com/jeonjonghyeok/coin/cli"
	"github.com/jeonjonghyeok/coin/db"
)

func main() {
	defer db.Close()
	cli.Start()

}
