package main

import (
	"fmt"
	"log"
	"os"

	"github.com/radean0909/redeam-rest/pkg/cmd"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	} else {
		log.Println("Starting Redeam Server")
	}
}
