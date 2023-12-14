package main

import (
	"flag"
	"fmt"

	"github.com/SashwatAnagolum/picodb/server"
)

func main() {
	var serverIPAddr, proxyIPAddr string

	flag.StringVar(&serverIPAddr, "server_ip",
		"127.0.0.1:4000", "IP Address and Port for the storage server")

	flag.StringVar(&proxyIPAddr, "proxy_ip",
		"127.0.0.1:3334", "IP Address and Port for the proxy server")

	flag.Parse()

	storageServer, err := server.NewStorageServer(
		serverIPAddr, proxyIPAddr,
		[]string{"requestcount", "requesttimes"},
		[]string{"disk"}, 100)

	if err != nil {
		fmt.Println(err)
		return
	}

	storageServer.Start()
}
