//+build zookeeper

package rpcxserver

import (
	"github.com/carefreeskyio/logger"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
)

func AddRegistryPlugin(s *server.Server, options *ServerOption) {
	r := &serverplugin.ZooKeeperRegisterPlugin{
		ServiceAddress:   options.Network + "@" + options.ServerIp + ":" + options.Port,
		ZooKeeperServers: options.RegistryAddr,
		BasePath:         options.BasePath,
		UpdateInterval:   options.UpdateInterval,
	}
	if err := r.Start(); err != nil {
		logger.Fatal(err)
	}
	s.Plugins.Add(r)
}
