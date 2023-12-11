package proxy

import (
	"fmt"
	"net"
	"sync"

	"github.com/SashwatAnagolum/picodb/utils"
	"google.golang.org/protobuf/proto"
)

type ClientRequest struct {
	clientConn *net.Conn
	request    *utils.PicoDBRequest
}

type ProxyServer struct {
	IPAddress      string
	Hash           *ConsistentHash
	Port           string
	ServerChannels map[uint32]chan []ClientRequest
	ChannelsLock   *sync.RWMutex
}

func (server *ProxyServer) cleanUpChannel(c chan []ClientRequest, storageServerID uint32) {
	server.ChannelsLock.Lock()
	delete(server.ServerChannels, storageServerID)
	server.ChannelsLock.Unlock()

	close(c)
}

func (server *ProxyServer) handleServerCommunication(
	conn *net.Conn) {
	storageServerID := server.Hash.hashServerIP(conn)

	storageServerCommChan := make(chan []ClientRequest, 256)

	server.ChannelsLock.Lock()
	server.Hash.AddToHashRing(storageServerID)
	server.ServerChannels[storageServerID] = storageServerCommChan
	server.ChannelsLock.Unlock()

	defer server.cleanUpChannel(storageServerCommChan, storageServerID)

	serverComm := NewServerCommunicator(conn, storageServerID, storageServerCommChan)
	// serverComm.Listen()
}

func (server *ProxyServer) serviceClientRequest(
	conn *net.Conn, buf []byte, len int) {
	// defer (*conn).Close()

	var request utils.PicoDBRequest
	proto.Unmarshal(buf[:len], &request)

	keyHash := server.Hash.hashString(request.Key)
	fmt.Println(server.Hash.GetNextServerIDOnRing(keyHash))

	responseBuffer, err := proto.Marshal(&utils.PicoDBResult{Key: "sammy", Value: "Boss"})

	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = (*conn).Write(responseBuffer)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func (server *ProxyServer) serviceRequest(conn *net.Conn) {
	buf := make([]byte, 2048)
	len, err := (*conn).Read(buf[:])

	if err != nil {
		fmt.Println(err.Error())
	}

	if len > 1 {
		server.serviceClientRequest(conn, buf, len)
	} else {
		server.handleServerCommunication(conn)
	}

}

func (server *ProxyServer) Listen() {
	listener, _ := net.Listen("tcp", server.IPAddress+":"+server.Port)

	for {
		connection, err := listener.Accept()

		fmt.Println("Received!")

		if err == nil {
			go server.serviceRequest(&connection)
		}
	}
}
