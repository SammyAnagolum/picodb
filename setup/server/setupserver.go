package main

import "github.com/SashwatAnagolum/picodb/server"

func main() {
	storageServer := server.NewStorageServer("127.0.0.1", "4000", "127.0.0.1:3333")
	storageServer.Listen()
}
