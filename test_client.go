package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/SashwatAnagolum/picodb/clientlib"
)

func main() {
	var proxyIPAddr string

	flag.StringVar(&proxyIPAddr, "proxy_ip",
		"127.0.0.1:3333", "IP Address and Port for the proxy server")

	dbClient, err := clientlib.NewPicoDBClient(proxyIPAddr)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(time.Now())

	for i := 0; i < 1; i++ {
		promise := dbClient.Put("SWENG421", "fun!")
		fmt.Println(promise.WaitForResult())
		promise = dbClient.Put("CMPSC445", "fun!")
		fmt.Println(promise.WaitForResult())
		promise = dbClient.Put("SWENG480", "fun!")
		fmt.Println(promise.WaitForResult())
		promise = dbClient.Put("CMPEN441", "fun!")
		fmt.Println(promise.WaitForResult())
	}

	fmt.Println(time.Now())
}
