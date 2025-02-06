package usecases

import (
	"fmt"
	"log/slog"
	"net"
	"pow-server/internal/infrastructure"
	"pow-server/internal/services"
)

type GetQuoteUsecase struct {
	Challenger    *services.Challenger
	QuoteProvider *services.QuoteProvider
	Logger        *slog.Logger
}

func NewGetQuoteUsecase(challenger *services.Challenger, quoteProvider *services.QuoteProvider, logger *slog.Logger) *GetQuoteUsecase {
	return &GetQuoteUsecase{
		Challenger:    challenger,
		QuoteProvider: quoteProvider,
		Logger:        logger,
	}
}

func (u *GetQuoteUsecase) Execute(conn net.Conn) error {

	//
	// 1. Challenge client
	//
	res, err := u.Challenger.Challenge(conn)
	if err != nil {
		return fmt.Errorf("challenging: %w", err)
	}

	u.Logger.Info("Challenge completed")

	if !res {
		return fmt.Errorf("chellange failed")
	}

	//
	// 2. Sending Payload
	//
	quote, err := u.QuoteProvider.GetQuote()
	if err != nil {
		return fmt.Errorf("get quote: %w", err)
	}

	//
	// 3. This one can be also extracted out here.
	//
	connAdapter := infrastructure.NewTcpAdapter()
	err = connAdapter.Write(conn, []byte(quote))
	if err != nil {
		return fmt.Errorf("sending quote: %w", err)
	}

	u.Logger.Info("Quote sent", "quote", quote)

	return nil
}
