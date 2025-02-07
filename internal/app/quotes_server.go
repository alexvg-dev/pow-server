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
		GetQuoteUsecase:    usecase,
		Cfg:                config,
		Logger:             logger,
		ConnectionsLimiter: make(chan struct{}, config.Server.MaxConnections),
	}
}

type QuotesApp struct {
	Logger             *slog.Logger
	Cfg                *config.Config
	TcpSrv             infrastructure.TcpServer
	GetQuoteUsecase    *usecases.GetQuoteUsecase
	ConnectionsLimiter chan struct{}
}

func (s *QuotesApp) Start(ctx context.Context) error {

	//
	// Starting metrics server (prometheus)
	//
	metricsServer := infrastructure.NewMetricsServer(s.Cfg.MetricsPort)
	s.Logger.Info("Starting metrics server", "port", s.Cfg.MetricsPort)
	go metricsServer.Start(ctx)

	//
	// Starting TCP Challenge server
	//
	tcpServer := infrastructure.NewTcpServer(s.Cfg.Server)
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
	// Handling new connections here
	//
	for {

		// It will be blocked until buffer if full
		s.ConnectionsLimiter <- struct{}{}

		conn, err := tcpServer.Accept()
		if errors.Is(err, net.ErrClosed) {
			s.Logger.Error("Listener closed", "err", err)
			return err
		}

		if err != nil {
			s.Logger.Error("Accept connection", "err", err)
		}

		s.Logger.Info("New connection", "from", conn.RemoteAddr().String())

		go func(curConn net.Conn) {
			defer curConn.Close()
			defer func() {
				<-s.ConnectionsLimiter
			}()

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
