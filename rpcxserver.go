package rpcxserver

import (
	"github.com/carefreeskyio/logger"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/util"
	"time"
)

type ServerOption struct {
	ServerName     string
	Network        string
	Port           string
	Registry       string
	BasePath       string
	UpdateInterval time.Duration
	RegistryAddr   []string
}

func NewServer(option *ServerOption) (s *server.Server) {
	addr, err := util.ExternalIPV4()
	if err != nil {
		logger.Fatalln("get ipv4 failed: err=%v", err)
	}

	s = server.NewServer()

	serverIp := addr + ":" + option.Port

	AddRegistryPlugin(s, option, serverIp)

	return s
}
