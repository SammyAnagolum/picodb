package main

import (
	"fmt"
	"net"

	"github.com/SashwatAnagolum/picodb/clientlib"
)

func main() {
	serverAddress, err := net.ResolveTCPAddr("tcp", "127.0.0.1:3333")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conn, err := net.DialTCP("tcp", nil, serverAddress)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dbClient := clientlib.NewPicoDBClient(conn)

	promise := dbClient.Put("SWENG421", "fun!")

	fmt.Println(promise.Ready())

	fmt.Println(promise.WaitForResult())

	// m := &utils.KeyValuePair{Key: "sammy", Value: "theboss"}
	// out, err := proto.Marshal(m)

	// if err == nil {
	// 	fmt.Println(out)
	// }

	// n := &utils.KeyValuePair{}
	// err = proto.Unmarshal(out, n)
	// fmt.Println(n)

	// concFilt1 := filters.EncryptKVPFilter{Data: &utils.KeyValuePair{Key: "sammy", Value: "puj"}}
	// concFilt2 := filters.EncryptKVPFilter{Source: &concFilt1}

	// out := concFilt2.GetData()

	// fmt.Println(out)
}
