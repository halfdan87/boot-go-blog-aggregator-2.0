package main

import (
	"fmt"

	"github.com/halfdan87/boot-go-blog-aggregator-2.0/internal/config"
)

func main() {
	fmt.Println("Starting...")

	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	cfg.SetUser("pioter")

	cfg, err = config.Read()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Current user: %v\n", cfg)

	fmt.Println("Done.")
}
