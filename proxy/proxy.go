package proxy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/SashwatAnagolum/picodb/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProxyServer struct {
	IPAddress          string
	ObserverIPAddress  string
	Hash               *ConsistentHash
	ServerClients      map[uint32]*utils.PicoDBServerClient
	HeartbeatChannels  map[uint32]chan bool
	ChannelsLock       *sync.RWMutex
	grpcServer         *grpc.Server
	grpcObserverServer *grpc.Server
	utils.UnimplementedPicoDBServerServer
	utils.UnimplementedPicoDBObserverServer
}

func NewProxyServer(address string, observerAddress string, numSlots uint32) *ProxyServer {
	ch := ConsistentHash{NumSlots: numSlots}

	server := ProxyServer{
		Hash:              &ch,
		IPAddress:         address,
		ObserverIPAddress: observerAddress,
		ServerClients:     make(map[uint32]*utils.PicoDBServerClient),
		HeartbeatChannels: make(map[uint32]chan bool),
		ChannelsLock:      &sync.RWMutex{}}

	return &server
}

func (server *ProxyServer) ProcessRequest(
	ctx context.Context,
	request *utils.PicoDBRequest) (*utils.PicoDBResponse, error) {
	keyHash := server.Hash.hashString(request.Key)
	serverID, err := server.Hash.GetNextServerIDOnRing(keyHash)

	if err != nil {
		fmt.Println(err)
		return &utils.PicoDBResponse{Key: request.Key}, err
	}

	server.ChannelsLock.RLock()
	storageServerClient, exists := server.ServerClients[serverID]
	server.ChannelsLock.RUnlock()

	if !exists {
		return &utils.PicoDBResponse{Key: request.Key}, errors.New(
			"storage server does not exist")
	}

	return (*storageServerClient).ProcessRequest(
		context.Background(), request)
}

func (server *ProxyServer) ManageHeartbeats(
	storageServerID uint32, hbChan chan bool) {
	shouldExit := false
	heartbeatTimer := time.NewTimer(10 * time.Second)

	for {
		select {
		case <-hbChan:
			heartbeatTimer.Stop()
			heartbeatTimer = time.NewTimer(10 * time.Second)
		case <-heartbeatTimer.C:
			server.ChannelsLock.Lock()
			delete(server.HeartbeatChannels, storageServerID)
			delete(server.ServerClients, storageServerID)
			server.ChannelsLock.Unlock()
			shouldExit = true
		}

		if shouldExit {
			break
		}
	}
}

func (server *ProxyServer) Notify(
	ctx context.Context,
	message *utils.PicoDBServerMessage) (*utils.PicoDBServerMessage, error) {
	storageServerID := server.Hash.hashString(message.SourceIPAddress)

	server.ChannelsLock.RLock()
	_, exists := server.HeartbeatChannels[storageServerID]
	server.ChannelsLock.RUnlock()

	if !exists {
		newServerClient, err := server.ConstructStorageServerClient(
			message.SourceIPAddress)

		if err != nil {
			fmt.Println(err)
			return &utils.PicoDBServerMessage{}, err
		} else {
			heartbeatChan := make(chan bool)

			server.ChannelsLock.Lock()
			server.ServerClients[storageServerID] = &newServerClient
			server.HeartbeatChannels[storageServerID] = heartbeatChan
			server.ChannelsLock.Unlock()

			server.Hash.AddToHashRing(storageServerID)
			go server.ManageHeartbeats(storageServerID, heartbeatChan)
		}
	}

	server.HeartbeatChannels[storageServerID] <- true

	return &utils.PicoDBServerMessage{}, nil
}

func (server *ProxyServer) ConstructStorageServerClient(
	storageServerAddress string) (utils.PicoDBServerClient, error) {
	var opts []grpc.DialOption
	opts = append(
		opts,
		grpc.WithTransportCredentials(
			insecure.NewCredentials()))

	conn, err := grpc.Dial(storageServerAddress, opts...)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	serverClient := utils.NewPicoDBServerClient(conn)

	return serverClient, nil
}

func (server *ProxyServer) StartObserverServer() {
	listener, err := net.Listen("tcp", server.ObserverIPAddress)

	if err != nil {
		fmt.Println(err)
		return
	}

	var serverOpts []grpc.ServerOption
	grpcServer := grpc.NewServer(serverOpts...)

	utils.RegisterPicoDBObserverServer(grpcServer, server)

	server.grpcObserverServer = grpcServer

	grpcServer.Serve(listener)
}

func (server *ProxyServer) StartServer() {
	listener, err := net.Listen("tcp", server.IPAddress)

	if err != nil {
		fmt.Println(err)
		return
	}

	var serverOpts []grpc.ServerOption
	grpcServer := grpc.NewServer(serverOpts...)

	utils.RegisterPicoDBServerServer(grpcServer, server)

	server.grpcServer = grpcServer

	grpcServer.Serve(listener)
}

func (server *ProxyServer) Start() {
	go server.StartObserverServer()
	server.StartServer()
}

func (server *ProxyServer) Close() {
	defer server.grpcServer.GracefulStop()
}
