package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"pow-server/internal/app"
	"pow-server/internal/config"
	"pow-server/internal/repository"
	"pow-server/internal/services"
	"pow-server/internal/usecases"
	"syscall"
)

func main() {

	if err := StarApp(); err != nil {
		panic(err)
	}
}

func StarApp() error {

	//
	// 1. Loading config
	//
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	//
	// 2. Init logger
	//
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logHandler)
	logger = logger.With("app", "Challenge server")
	logger.Info("Starting app", "port", cfg.Server.Port)

	//
	// 3. Building server
	//
	challenger := services.NewChallenger()
	quoteRepo := repository.NewQuotesRepo(cfg.QuotesFilePath)
	quoteProvider := services.NewQuoteProvider(quoteRepo)
	usecase := usecases.NewGetQuoteUsecase(challenger, quoteProvider, logger)

	//
	// 4. Gracefull shutdown
	//
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	//
	// 5. Starting server
	//
	quotesApp := app.NewQuotesApp(cfg, usecase, logger)
	err = quotesApp.Start(ctx)
	if err != nil {
		logger.Error("Starting server", "error", err)
		return err
	}

	return nil
}
