package clientlib

import (
	"github.com/SashwatAnagolum/picodb/clientlib/filters"
	"github.com/SashwatAnagolum/picodb/clientlib/promise"
	"github.com/SashwatAnagolum/picodb/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PicoDBClient struct {
	ServerAddress    string
	RequestFilter    filters.KVPFilterIF
	ServerConnection *grpc.ClientConn
	ServerClient     *utils.PicoDBServerClient
}

func NewPicoDBClient(serverAddress string) (*PicoDBClient, error) {
	var opts []grpc.DialOption
	opts = append(
		opts,
		grpc.WithTransportCredentials(
			insecure.NewCredentials()))

	conn, err := grpc.Dial(serverAddress, opts...)

	if err != nil {
		return nil, err
	}

	serverClient := utils.NewPicoDBServerClient(conn)

	client := &PicoDBClient{
		ServerAddress: serverAddress,
		ServerClient:  &serverClient,
		RequestFilter: &filters.IdentityKVPFilter{
			AbsKVPFilter: filters.AbsKVPFilter{}}}

	return client, nil
}

func (client *PicoDBClient) CloseServerConnection() {
	defer client.ServerConnection.Close()
}

func (client *PicoDBClient) TransformData(
	kvp *utils.PicoDBRequest) *utils.PicoDBRequest {
	client.RequestFilter.SetData(kvp)

	return client.RequestFilter.GetData()
}

func (client *PicoDBClient) sendRequest(
	request *utils.PicoDBRequest) *promise.PicoDBResponsePromise {
	promise := promise.NewResultPromise()
	promise.Start(client.ServerClient, request)

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
	key string, value string) *promise.PicoDBResponsePromise {
	request := client.makeRequest(key, value, utils.PUT)
	promise := client.sendRequest(request)

	return promise
}

func (client *PicoDBClient) Get(
	key string) *promise.PicoDBResponsePromise {
	request := client.makeRequest(key, "", utils.GET)
	promise := client.sendRequest(request)

	return promise
}

func (client *PicoDBClient) Delete(
	key string) *promise.PicoDBResponsePromise {
	request := client.makeRequest(key, "", utils.DELETE)
	promise := client.sendRequest(request)

	return promise
}
