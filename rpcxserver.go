package rpcxserver

import (
	"github.com/smallnest/rpcx/server"
	"time"
)

type ServerOption struct {
	ServerName     string
	ServerIp       string
	Network        string
	Port           string
	BasePath       string
	UpdateInterval time.Duration
	RegistryAddr   []string
}

func NewServer(option *ServerOption) (s *server.Server) {
	s = server.NewServer()

	AddRegistryPlugin(s, option)

	return s
}
