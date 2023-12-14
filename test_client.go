package main

import (
	"flag"
	"fmt"

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

	promise := dbClient.Put("SWENG421", "5*")
	fmt.Println(promise.WaitForResult())

	promise = dbClient.Put("CMPSC445", "5*")
	fmt.Println(promise.WaitForResult())

	promise = dbClient.Put("SC120N", "5*")
	fmt.Println(promise.WaitForResult())

	promise = dbClient.Put("EARTH104N", "5*")
	fmt.Println(promise.WaitForResult())

	promise = dbClient.Put("COMM150N", "5*")
	fmt.Println(promise.WaitForResult())

	promise = dbClient.Put("CMPSC122", "5*")
	fmt.Println(promise.WaitForResult())

	promise = dbClient.Put("SWENG480", "5*")
	fmt.Println(promise.WaitForResult())
}
