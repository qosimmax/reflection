package main

import (
	"flag"
	"fmt"
	"os"
	"reflection/parser"
)

func main() {
	arg := &parser.Args{}

	err := parser.GetArguments(arg)
	if err == parser.ErrRequired {
		fmt.Println("flag parse error", err)
		flag.PrintDefaults()
		os.Exit(1)
		return
	}

	fmt.Printf("%#v", arg)
}
