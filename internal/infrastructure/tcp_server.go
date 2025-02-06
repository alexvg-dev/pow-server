package infrastructure

import (
	"context"
	"net"
	"pow-server/internal/config"
	"strconv"
)

type TcpServer struct {
	Settings config.ServerConfig
	Listener net.Listener
}

func NewTcpServer(cfg config.ServerConfig) *TcpServer {
	return &TcpServer{
		Settings: cfg,
	}
}

func (srv *TcpServer) Start(runCtx context.Context) error {

	// Keep-alive settings can be places here
	// 	-- idle
	//	-- internal
	//	-- count
	//
	listenCfg := net.ListenConfig{
		KeepAlive: 0,
	}
	listener, err := listenCfg.Listen(runCtx, "tcp", ":"+strconv.Itoa(srv.Settings.Port))
	if err != nil {
		return err
	}

	srv.Listener = listener

	return nil
}

func (srv *TcpServer) Accept() (net.Conn, error) {
	return srv.Listener.Accept()
}

func (srv *TcpServer) Stop() {
	if nil != srv.Listener {
		srv.Listener.Close()
	}
}
