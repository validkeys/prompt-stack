package main

import (
	"flag"
	"fmt"
	"os"

	"example.com/style-anchor/pkg/greeter"
)

func main() {
	name := flag.String("name", "world", "name to greet")
	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Println(greeter.Hello(*name))
}
