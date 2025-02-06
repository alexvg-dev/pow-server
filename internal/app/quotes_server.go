package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"pow-server/internal/config"
	"pow-server/internal/infrastructure"
	"pow-server/internal/usecases"
	"time"
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

	// TODO: add metrics server
	tcpServer := infrastructure.NewTcpServer(s.Cfg.Server)

	//
	// Starting port listening
	//
	err := tcpServer.Start(ctx)
	if err != nil {
		s.Logger.Error("Waiting for new connections", "err", err)
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
		// TODO: limit max connections by config
		conn, err := tcpServer.Accept()
		if errors.Is(err, net.ErrClosed) {
			s.Logger.Error("Listener closed", "err", err)
			return nil
		}

		if err != nil {
			s.Logger.Error("Accept connection", "err", err)
		}

		s.Logger.Info("New connection", "from", conn.RemoteAddr().String())

		go func(curConn net.Conn) {
			defer curConn.Close()

			// It will help us to terminate usecase execution if TTL expired
			//
			ttlCtx, cancel := context.WithTimeout(context.Background(), time.Duration(s.Cfg.SessionTtlSec)*time.Second)
			defer cancel()

			err := s.GetQuoteUsecase.Execute(ttlCtx, curConn)
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
