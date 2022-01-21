package rpcxserver

import (
	"context"
	"fmt"
	"github.com/carefreeskyio/logger"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
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
	Name    string
	Addr    string
	Network string
	Port    string
}

type RegistryOption struct {
	Addr           []string
	BasePath       string
	UpdateInterval time.Duration
	Group          string
}

type RateLimitOption struct {
	Enable       bool
	FillInterval time.Duration
	Capacity     int64
}

type Options struct {
	Server    ServerOption
	Registry  RegistryOption
	Service   interface{}
	RateLimit RateLimitOption
	Plugin    []server.Plugin
}

func NewServer(options *Options) *RpcXServer {
	s := server.NewServer()

	AddRegistryPlugin(s, options)

	if options.RateLimit.Enable {
		options.Plugin = append(options.Plugin, serverplugin.NewReqRateLimitingPlugin(options.RateLimit.FillInterval, options.RateLimit.Capacity, true))
	}

	addPlugins(s, options.Plugin)

	if err := s.RegisterName(options.Server.Name, options.Service, "group="+options.Registry.Group); err != nil {
		logger.Fatalf("start service failed: err=%v", err)
	}

	return &RpcXServer{
		ServerName: options.Server.Name,
		Server:     s,
	}
}

func addPlugins(server *server.Server, plugins []server.Plugin) {
	if len(plugins) == 0 {
		return
	}
	for _, p := range plugins {
		server.Plugins.Add(p)
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
