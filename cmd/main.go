package main

import (
	"LOTestTask/internal/di"
	"context"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	di := di.DI{}
	if err := di.Init(); err != nil {
		log.Fatalf("can not init service - %v", err)
	}

	go func() {
		if err := di.Start(); err != nil {
			log.Printf("error ocured while starting server - %v", err)
			cancel()
		}
	}()

	di.Stop(ctx)
}
