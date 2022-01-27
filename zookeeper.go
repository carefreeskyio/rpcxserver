//+build zookeeper

package rpcxserver

import (
	"github.com/carefreex-io/logger"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"time"
)

func AddRegistryPlugin(s *server.Server, options **BaseOptions) {
	r := &serverplugin.ZooKeeperRegisterPlugin{
		ServiceAddress:   options.Server.Network + "@" + options.Server.Addr + ":" + options.Server.Port,
		ZooKeeperServers: options.Registry.Addr,
		BasePath:         options.Registry.BasePath,
		UpdateInterval:   options.Registry.UpdateInterval * time.Second,
	}
	if err := r.Start(); err != nil {
		logger.Fatal(err)
	}
	s.Plugins.Add(r)
}
