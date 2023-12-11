package clientlib

import (
	"net"

	"github.com/SashwatAnagolum/picodb/clientlib/filters"
	"github.com/SashwatAnagolum/picodb/clientlib/promise"
	"github.com/SashwatAnagolum/picodb/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PicoDBClient struct {
	ServerConnection *net.TCPConn
	RequestFilter    filters.KVPFilterIF
}

func NewPicoDBClient(connection *net.TCPConn) PicoDBClient {
	client := PicoDBClient{
		ServerConnection: connection,
		RequestFilter: &filters.IdentityKVPFilter{
			AbsKVPFilter: filters.AbsKVPFilter{}}}

	return client
}

func (client *PicoDBClient) TransformData(
	kvp *utils.PicoDBRequest) *utils.PicoDBRequest {
	client.RequestFilter.SetData(kvp)

	return client.RequestFilter.GetData()
}

func (client *PicoDBClient) sendRequest(
	request *utils.PicoDBRequest) *promise.PicoDBResultPromise {
	promise := promise.NewResultPromise()
	promise.Start(client.ServerConnection, request)

	return promise
}

func (client *PicoDBClient) makeRequest(
	key string, value string,
	requestType utils.PicoDBRequestType) *utils.PicoDBRequest {
	request := utils.PicoDBRequest{
		Key:         key,
		Value:       value,
		Time:        timestamppb.Now(),
		RequestType: requestType.EnumIndex()}

	return client.TransformData(&request)
}

func (client *PicoDBClient) Put(
	key string, value string) *promise.PicoDBResultPromise {
	request := client.makeRequest(key, value, utils.PUT)
	promise := client.sendRequest(request)

	return promise
}

func (client *PicoDBClient) Get(
	key string) *promise.PicoDBResultPromise {
	request := client.makeRequest(key, "", utils.GET)
	promise := client.sendRequest(request)

	return promise
}

func (client *PicoDBClient) Delete(
	key string) *promise.PicoDBResultPromise {
	request := client.makeRequest(key, "", utils.DELETE)
	promise := client.sendRequest(request)

	return promise
}
