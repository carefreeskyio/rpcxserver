//+build etcd

package rpcxserver

import (
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
)

func AddRegistryPlugin(s *server.Server, options *ServerOption, serverIp string) {
	r := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: options.Network + "@" + serverIp,
		EtcdServers:    options.RegistryAddr,
		BasePath:       options.BasePath,
		UpdateInterval: options.UpdateInterval,
	}
	if err := r.Start(); err != nil {
		logger.Fatal(err)
	}
	s.Plugins.Add(r)
}
