package main

import (
	"fmt"
	"log/slog"
	"os"
	"pow-server/pkg/pow"
	"pow-server/pkg/quotes_client"
	"sync"
)

const (
	PoWDifficulty = 2 // Количество нулей в начале хэша
	ClientsNumber = 8
)

func main() {

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	serverAddress := os.Args[1]

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logHandler)
	logger = logger.With("app", "Challenge client")

	wg := sync.WaitGroup{}
	wg.Add(ClientsNumber)

	//
	// Starting multiple clients to check possible
	//	concurrency problems
	//
	for i := 0; i < ClientsNumber; i++ {
		go func(clientID int) {
			defer wg.Done()

			logger = logger.With("ClientNum", clientID)
			logger.Info("Starting client", "num", clientID)

			powSolver := pow.NewScryptPow(PoWDifficulty)
			client := quotes_client.NewClient(serverAddress, powSolver)

			logger.Info("Fetching quote")

			quote, err := client.GetQuote()
			if err != nil {
				logger.Error("Get quote", "err", err)
				return
			}

			logger.Info("Got quote", "quote", quote)
		}(i)
	}

	wg.Wait()
}

func printUsage() {
	fmt.Println("")
	fmt.Println("Usage: ./client <server_address:port>")
	fmt.Println("Example: ./client localhost:4444")
	fmt.Println("")
}
