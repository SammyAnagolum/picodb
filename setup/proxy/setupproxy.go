package main

import (
	"flag"

	"github.com/SashwatAnagolum/picodb/proxy"
)

func main() {
	var serverIPAddr, observerIPAddr string

	flag.StringVar(&serverIPAddr, "proxy_ip",
		"127.0.0.1:3333", "IP Address and Port for the Proxy server")

	flag.StringVar(&observerIPAddr, "proxy_observer_ip",
		"127.0.0.1:3334", "IP Address and Port for the storage node observer server")

	flag.Parse()

	proxyServer := proxy.NewProxyServer(serverIPAddr, observerIPAddr, 1024)
	proxyServer.Start()
}
