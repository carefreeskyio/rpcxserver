//+build etcd

package rpcxserver

import (
	"github.com/carefreeskyio/logger"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"time"
)

func AddRegistryPlugin(s *server.Server, options *ServerOption) {
	r := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: options.Network + "@" + options.ServerIp + ":" + options.Port,
		EtcdServers:    options.RegistryAddr,
		BasePath:       options.BasePath,
		UpdateInterval: options.UpdateInterval * time.Minute,
	}
	if err := r.Start(); err != nil {
		logger.Fatal(err)
	}
	s.Plugins.Add(r)
}
