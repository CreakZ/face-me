package main

import (
	"context"
	"faceit_parser/parser"
	"fmt"
	"time"
)

func main() {
	start := time.Now()

	var nickname = "s1mple"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	matches, err := parser.GetAllMatchesStats(ctx, nickname, "csgo")
	if err != nil {
		panic(fmt.Errorf("%v", err))
	}

	fmt.Println(len(matches))

	fmt.Printf("exec time: %v\n", time.Since(start))
}
