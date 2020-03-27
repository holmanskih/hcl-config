package main

import (
	"log"

	"github.com/holmanskih/hcl-config/parse"
)

func main() {
	cfg, err := parse.LoadConfig("env")
	if err != nil {
		log.Fatalf("failed file parsing: %s", err)
	}

	log.Printf("loaded config %v", cfg)
}
