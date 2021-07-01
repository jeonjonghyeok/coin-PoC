package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/jeonjonghyeok/coin/rest"
)

//github.com/spf13/cobra

func usage() {
	fmt.Printf("Welcome to Coin\n")
	fmt.Printf("Please use the follwing flags:\n")
	fmt.Printf("-port:		Set the PORT of the server\n")
	fmt.Printf("-mode:		Start the REST API (recommended)\n\n")
	runtime.Goexit()
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}
	port := flag.Int("port", 5000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")
	flag.Parse()
	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
	//	explorer.Start(*port)
	default:
		usage()
	}

}
