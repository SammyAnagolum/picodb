package main

import (
	"sync"

	"github.com/SashwatAnagolum/picodb/proxy"
)

func main() {
	consistentHash := proxy.ConsistentHash{NumSlots: 1024}
	server := proxy.ProxyServer{
		Hash:           &consistentHash,
		IPAddress:      "127.0.0.1",
		Port:           "3333",
		ServerChannels: make(map[uint32]chan []proxy.ClientRequest),
		ChannelsLock:   &sync.RWMutex{}}

	server.Listen()
}
