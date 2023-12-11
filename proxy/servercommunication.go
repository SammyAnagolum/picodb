package proxy

import (
	"fmt"
	"net"
)

type ServerCommunicator struct {
	ServerConn *net.Conn
	ServerID   uint32
	CRChan     chan []ClientRequest
}

func NewServerCommunicator(serverConn *net.Conn,
	serverID uint32, crChan chan []ClientRequest) *ServerCommunicator {
	return &ServerCommunicator{
		ServerConn: serverConn, ServerID: serverID, CRChan: crChan}
}

func (comm *ServerCommunicator) Notify() {

}

func (comm *ServerCommunicator) WriteToStorageServer() {

}

func (comm *ServerCommunicator) ReadFromStorageServer() {
	buf := make([]byte, 2048)

	for {
		len, err := (*comm.ServerConn).Read(buf)
		fmt.Println(buf[:len], err)

		if err != nil {
			break
		}

		if len == 1 {
			comm.Notify()
		}
	}
}

func (comm *ServerCommunicator) Start() {

}
