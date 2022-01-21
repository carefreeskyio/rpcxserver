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

	options := Options{
		Server: ServerOption{
			Name:    "test",
			Addr:    addr,
			Network: "tcp",
			Port:    "8090",
		},
		Registry: RegistryOption{
			BasePath:       "/test",
			UpdateInterval: 60,
			Addr:           []string{"127.0.0.1:2379"},
			Group:          "test",
		},
		Service: &Service{},
	}
	s := NewServer(&options)

	s.Start(options.Server.Network, addr+":"+options.Server.Port)
}
