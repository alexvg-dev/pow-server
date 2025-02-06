package usecases

import (
	"context"
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

func (u *GetQuoteUsecase) Execute(ctx context.Context, conn net.Conn) error {

	//
	// We have two options of processing timeouts here:
	//	1. start goroutine which waits for ctx.Done() and closes connection
	//		- it stop execution with error and may have unpredictable behavior
	//	2. check context after each step manually and stop execution explicitly -
	//		- which is better than first one, but more wordily
	//

	//
	// 1. Challenge client
	//
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("timeout befor challange: %w", err)
	}

	res, err := u.Challenger.Challenge(ctx, conn)
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
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("timeout befor sending payload: %w", err)
	}
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
