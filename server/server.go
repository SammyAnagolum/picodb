package server

import (
	"fmt"
	"net"
	"time"

	"github.com/SashwatAnagolum/picodb/utils"
	"google.golang.org/protobuf/proto"
)

type PicoDBStorageServer struct {
	Data            map[string]string
	ProxyServerConn *net.TCPConn
	IPAddress       string
	Port            string
}

func NewStorageServer(ipAddress string, port string, proxyServerAddress string) *PicoDBStorageServer {
	proxyServerAddr, err := net.ResolveTCPAddr("tcp", proxyServerAddress)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	conn, err := net.DialTCP("tcp", nil, proxyServerAddr)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return &PicoDBStorageServer{
		IPAddress:       ipAddress,
		Port:            port,
		Data:            map[string]string{},
		ProxyServerConn: conn}
}

func (server *PicoDBStorageServer) serviceRequest(conn *net.Conn) {
	defer (*conn).Close()

	buf := make([]byte, 2048)
	len, err := (*conn).Read(buf)

	if err != nil {
		(*conn).Write([]byte("An error occured! Please try again."))
	}

	var request utils.PicoDBRequest
	proto.Unmarshal(buf[:len], &request)

	fmt.Println(request)
}

func (server *PicoDBStorageServer) sendHeartbeats() {
	byteArr := make([]byte, 1)
	byteArr[0] = 0xFF

	for {
		fmt.Println("Sending!")
		_, err := server.ProxyServerConn.Write(byteArr)

		if err != nil {
			fmt.Println(err.Error())
		}

		time.Sleep(2 * time.Second)
	}
}

func (server *PicoDBStorageServer) Listen() {
	go server.sendHeartbeats()

	listener, _ := net.Listen("tcp", server.IPAddress+":"+server.Port)

	for {
		connection, err := listener.Accept()

		if err == nil {
			go server.serviceRequest(&connection)
		}
	}
}
