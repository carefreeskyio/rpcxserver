package rpcxserver

import (
	"context"
	"fmt"
	"github.com/carefreeskyio/logger"
	"github.com/smallnest/rpcx/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RpcXServer struct {
	ServerName       string
	Server           *server.Server
	onStartAction    []func(s *server.Server)
	onShutdownAction []func(s *server.Server)
}

type ServerOption struct {
	ServerName     string
	ServerIp       string
	Network        string
	Port           string
	BasePath       string
	UpdateInterval time.Duration
	Group          string
	RegistryAddr   []string
	Service        interface{}
}

func NewServer(option *ServerOption) *RpcXServer {

	s := server.NewServer()

	AddRegistryPlugin(s, option)

	if err := s.RegisterName(option.ServerName, option.Service, "group="+option.Group); err != nil {
		logger.Fatalf("start service failed: err=%v", err)
	}

	return &RpcXServer{
		ServerName: option.ServerName,
		Server:     s,
	}
}

func (s *RpcXServer) AddOnStartAction(fn func(s *server.Server)) {
	s.onStartAction = append(s.onStartAction, fn)
}

func (s *RpcXServer) AddOnShutdownAction(fn func(s *server.Server)) {
	s.onShutdownAction = append(s.onShutdownAction, fn)
}

func (s *RpcXServer) Start(network string, address string) {
	s.onStart()

	go func() {
		if err := s.Server.Serve(network, address); err != nil {
			if err == server.ErrServerClosed {
				logger.Info(err)
			} else {
				panic(err)
			}
		}
	}()
	fmt.Println(s.ServerName + " start successfully")

	s.waitShutdown()
}

func (s *RpcXServer) onStart() {
	for _, fn := range s.onStartAction {
		fn(s.Server)
	}
}

func (s *RpcXServer) onShutdown() {
	for _, fn := range s.onShutdownAction {
		fn(s.Server)
	}
}

func (s *RpcXServer) waitShutdown() {
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	<-sig

	s.onShutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := s.Server.Shutdown(ctx); err != nil {
		logger.Errorf("call s.Server.Shutdown failed: err=%v", err)
	}
	cancel()
}
