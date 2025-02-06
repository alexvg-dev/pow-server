package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"pow-server/internal/config"
	"pow-server/internal/infrastructure"
	"pow-server/internal/usecases"
)

func NewQuotesApp(config *config.Config, usecase *usecases.GetQuoteUsecase, logger *slog.Logger) *QuotesApp {
	return &QuotesApp{
		GetQuoteUsecase: usecase,
		Cfg:             config,
		Logger:          logger,
	}
}

type QuotesApp struct {
	Logger          *slog.Logger
	Cfg             *config.Config
	TcpSrv          infrastructure.TcpServer
	GetQuoteUsecase *usecases.GetQuoteUsecase
}

func (s *QuotesApp) Start(ctx context.Context) error {

	tcpServer := infrastructure.NewTcpServer(s.Cfg.Server)

	//
	// Starting port listening
	//
	err := tcpServer.Start(ctx)
	if err != nil {
		s.Logger.Error("Waiting for new connections", err)
	}

	go func() {
		select {
		case <-ctx.Done():
			s.Logger.Info("Stopping server")
			tcpServer.Stop()
		}
	}()

	//
	// Handling connections here
	//
	for {
		conn, err := tcpServer.Accept()
		if errors.Is(err, net.ErrClosed) {
			s.Logger.Error("Listener closed", "err", err)
			return nil
		}

		if err != nil {
			s.Logger.Error("Accept connection", err)
		}

		s.Logger.Info("New connection", "from", conn.RemoteAddr().String())

		go func(curConn net.Conn) {
			defer curConn.Close()

			err := s.GetQuoteUsecase.Execute(curConn)
			if err != nil {
				s.Logger.Error("Get quote usecase execute", "err", err)
				return
			}

			return
		}(conn)
	}
}

func (s *QuotesApp) Stop() error {
	s.TcpSrv.Stop()

	return nil
}
