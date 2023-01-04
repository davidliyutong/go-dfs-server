package server

import (
	"go-dfs-server/pkg/config"
	v1 "go-dfs-server/pkg/dataserver/client/v1"
	"go-dfs-server/pkg/status"
)

type dataServerManager struct {
	desc         *config.NameServerDesc
	clients      []interface{}
	namedClients map[string]interface{}
}

func (d *dataServerManager) GetAllClients() []interface{} {
	return d.clients
}

func (d *dataServerManager) GetClient(uuid string) (interface{}, error) {
	client, ok := d.namedClients[uuid]
	if !ok {
		return nil, status.ErrClientNotFound
	} else {
		return client, nil
	}
}

func (d *dataServerManager) GetClients(uuids []string) ([]interface{}, error) {
	clients := make([]interface{}, 0, 3)
	var errOccurred = false
	for _, uuid := range uuids {
		client, err := d.GetClient(uuid)
		if err != nil {
			errOccurred = true
		} else {
			clients = append(clients, client)
		}
	}
	if errOccurred {
		return clients, status.ErrClientNotFoundSome
	} else {
		return clients, nil
	}
}

func (d *dataServerManager) ListServers() []config.RegisteredDataServer {
	return d.desc.Opt.DataServers
}

func (d *dataServerManager) AlivenessProbe() (map[string]bool, error) {
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
		return stat, status.ErrDataServerOfflineSome
	}
}

func (d *dataServerManager) UUIDProbe() (map[string]bool, error) {
	stat := make(map[string]bool, len(d.clients))
	var flag = true
	for idx, client := range d.clients {
		if client == nil {
			continue
		}
		serverUUID, err := client.(v1.DataServerClient).SysUUID()
		if err != nil {
			flag = false
			continue
		}
		if serverUUID == "" || serverUUID != client.(v1.DataServerClient).GetUUID() && client.(v1.DataServerClient).GetUUID() != "" {
			flag = false
		} else {
			if client.(v1.DataServerClient).GetUUID() == "" {
				stat[serverUUID] = true
				client.(v1.DataServerClient).SetUUID(serverUUID)
				d.desc.Opt.DataServers[idx].UUID = serverUUID
				d.namedClients[serverUUID] = client
			}
		}
	}
	if flag {
		return stat, nil
	} else {
		return stat, status.ErrDataServerUUIDMismatch
	}
}

func (d *dataServerManager) Register() error {

	for _, client := range d.clients {
		if client == nil {
			continue
		}
		olderUUID, err := client.(v1.DataServerClient).SysRegister(GlobalServerDesc.UUID)
		if err != nil {
			return err
		}
		if olderUUID == "" {
			return status.ErrDataServerReboot
		} else if olderUUID != GlobalServerDesc.UUID {
			return status.ErrNameServerReboot
		}
	}
	return nil

}

type DataServerManager interface {
	ListServers() []config.RegisteredDataServer
	AlivenessProbe() (map[string]bool, error)
	UUIDProbe() (map[string]bool, error)
	Register() error
	GetAllClients() []interface{}
	GetClient(uuid string) (interface{}, error)
	GetClients(uuid []string) ([]interface{}, error)
}

func NewDataServerManager(desc *config.NameServerDesc) DataServerManager {
	clients := make([]interface{}, 0)
	namedClients := make(map[string]interface{})
	for _, server := range desc.Opt.DataServers {
		newClient := v1.NewDataServerClient(server.UUID, server.Address, server.Port, false)
		clients = append(clients, newClient)
		if server.UUID != "" {
			namedClients[server.UUID] = newClient
		}
	}
	return &dataServerManager{desc: desc, clients: clients, namedClients: namedClients}
}
