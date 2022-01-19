package rpcxserver

import (
	"context"
	"github.com/smallnest/rpcx/util"
	"testing"
)

type Service struct {
}

func (s *Service) Hello(ctx context.Context, request *int, response *int) error {
	return nil
}

func TestNewServer(t *testing.T) {
	addr, _ := util.ExternalIPV4()

	serverOption := ServerOption{
		ServerName:     "test",
		ServerIp:       addr,
		Network:        "tcp",
		Port:           "8090",
		BasePath:       "/test",
		UpdateInterval: 1,
		RegistryAddr:   []string{"127.0.0.1:2379"},
		Group:          "test",
		Service:        &Service{},
	}
	s := NewServer(&serverOption)

	s.Start(serverOption.Network, addr+":"+serverOption.Port)
}
