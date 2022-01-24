package rpcxserver

import (
	"context"
	"fmt"
	"github.com/carefreex-io/config"
	"github.com/carefreex-io/logger"
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

type BaseOptions struct {
	Server    ServerOption
	Registry  RegistryOption
	RateLimit RateLimitOption
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

type CustomOptions struct {
	Service interface{}
	Plugin  []server.Plugin
}

var (
	baseOptions          *BaseOptions
	DefaultCustomOptions = &CustomOptions{}
)

func initBaseOptions() {
	baseOptions = &BaseOptions{
		Server: ServerOption{
			Name:    config.GetString("Service.Name"),
			Network: config.GetString("Service.Network"),
			Port:    config.GetString("Service.Port"),
		},
		Registry: RegistryOption{
			Addr:           config.GetStringSlice("Registry.Addr"),
			BasePath:       config.GetString("Registry.BasePath"),
			UpdateInterval: config.GetDuration("Registry.UpdateInterval"),
			Group:          config.GetString("Registry.Group"),
		},
		RateLimit: RateLimitOption{
			Enable:       config.GetBool("RateLimit.Enable"),
			FillInterval: config.GetDuration("RateLimit.FillInterval"),
			Capacity:     config.GetInt64("RateLimit.Token"),
		},
	}
}

func NewServer() *RpcXServer {
	initBaseOptions()

	s := server.NewServer()

	AddRegistryPlugin(s, baseOptions)

	if config.GetBool("RateLimit.Enable") {
		DefaultCustomOptions.Plugin = append(DefaultCustomOptions.Plugin, serverplugin.NewReqRateLimitingPlugin(baseOptions.RateLimit.FillInterval, baseOptions.RateLimit.Capacity, true))
	}

	addPlugins(s, DefaultCustomOptions.Plugin)

	if err := s.RegisterName(baseOptions.Server.Name, DefaultCustomOptions.Service, "group="+baseOptions.Registry.Group); err != nil {
		logger.Fatalf("start service failed: err=%v", err)
	}

	return &RpcXServer{
		ServerName: baseOptions.Server.Name,
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
