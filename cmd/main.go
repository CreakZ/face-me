package main

import (
	"faceit_parser/internal/delivery"
	"log"
)

func main() {
	logger := log.Default()
	// info log

	delivery.Start(logger)
}
