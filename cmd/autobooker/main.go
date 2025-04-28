package main

import (
	"log"

	"github.com/ekkx/tcmrsv/cmd/autobooker/commands"
)

func main() {
	if err := commands.Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
