package server

import (
	"errors"
	"go-dfs-server/pkg/config"
	v1 "go-dfs-server/pkg/dataserver/client/v1"
)

type dataServerManager struct {
	desc    *config.NameServerDesc
	clients []interface{}
}

func (d *dataServerManager) ListServers() []config.RegisteredDataServer {
	//TODO implement me
	return d.desc.Opt.DataServers
}

func (d *dataServerManager) AlivenessProbe() (map[string]bool, error) {
	//TODO implement me
	stat := make(map[string]bool, len(d.clients))
	var flag bool
	for _, client := range d.clients {
		err := client.(v1.DataServerClient).Ping()
		if err != nil {
			flag = false
		}
		stat[client.(v1.DataServerClient).GetUUID()] = err == nil
	}
	if flag {
		return stat, nil
	} else {
		return stat, errors.New("not all data servers are alive")
	}
}

func (d *dataServerManager) UUIDProbe() (map[string]bool, error) {
	//TODO implement me
	stat := make(map[string]bool, len(d.clients))
	var flag bool
	for _, client := range d.clients {
		serverUUID, err := client.(v1.DataServerClient).SysUUID()
		if err != nil {
			flag = false
			stat[client.(v1.DataServerClient).GetUUID()] = false
			continue
		}
		if serverUUID != client.(v1.DataServerClient).GetUUID() {
			flag = false
			stat[client.(v1.DataServerClient).GetUUID()] = false
		} else {
			stat[client.(v1.DataServerClient).GetUUID()] = true
		}
	}
	if flag {
		return stat, nil
	} else {
		return stat, errors.New("not all data servers uuid match record")
	}
}

type DataServerManager interface {
	ListServers() []config.RegisteredDataServer
	AlivenessProbe() (map[string]bool, error)
	UUIDProbe() (map[string]bool, error)
}

func NewDataServerManager(desc *config.NameServerDesc) DataServerManager {
	clients := make([]interface{}, len(desc.Opt.DataServers))
	for _, server := range desc.Opt.DataServers {
		clients = append(clients, v1.NewDataServerClient(server.UUID, server.Address, server.Port, false))
	}
	return &dataServerManager{desc: desc, clients: clients}
}
