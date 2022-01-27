//+build etcdV3

package rpcxserver

import (
	"github.com/carefreex-io/logger"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"time"
)

func AddRegistryPlugin(s *server.Server, options *BaseOptions) {
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: options.Server.Network + "@" + options.Server.Addr + ":" + options.Server.Port,
		EtcdServers:    options.Registry.Addr,
		BasePath:       options.Registry.BasePath,
		UpdateInterval: options.Registry.UpdateInterval * time.Second,
	}
	if err := r.Start(); err != nil {
		logger.Fatal(err)
	}
	s.Plugins.Add(r)
}
