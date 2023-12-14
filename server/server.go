package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/SashwatAnagolum/picodb/server/logging"
	"github.com/SashwatAnagolum/picodb/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PicoDBStorageServer struct {
	Data            map[string]string
	ProxyServerAddr string
	IPAddress       string
	Logger          logging.LoggerIF
	grpcServer      *grpc.Server
	proxyClient     utils.PicoDBObserverClient
	utils.UnimplementedPicoDBServerServer
}

func NewStorageServer(address string,
	proxyServerAddress string, loggers []string,
	writers []string, writeFreq int) (*PicoDBStorageServer, error) {
	logger, err := getLogger(loggers)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	logWriters, err := getLogWriter(writers)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, logWriter := range logWriters {
		logger = logging.NewLoggerWriterWrapper(logger, logWriter, writeFreq)
	}

	return &PicoDBStorageServer{
		IPAddress:       address,
		Data:            map[string]string{},
		ProxyServerAddr: proxyServerAddress,
		Logger:          logger}, nil
}

func getLogger(loggers []string) (logging.LoggerIF, error) {
	loggerFactory := logging.NewLoggerFactory()
	var logger logging.LoggerIF
	var err error

	if len(loggers) > 0 {
		logger, err = loggerFactory.NewLogger(loggers[0])

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for _, nextLoggerName := range loggers[1:] {
			nextLogger, err := loggerFactory.NewLogger(nextLoggerName)

			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			logger = logging.NewLoggerWrapper(logger, nextLogger)
		}
	} else {
		logger, err = loggerFactory.NewLogger("null")

		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return logger, nil
}

func getLogWriter(writers []string) ([]logging.LogWriterIF, error) {
	writerList := make([]logging.LogWriterIF, 0)
	writerFactory := logging.NewLogWriterFactory()

	for _, writerName := range writers {
		newWriter, err := writerFactory.NewLogWriter(writerName)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		writerList = append(writerList, newWriter)
	}

	return writerList, nil
}

func (server *PicoDBStorageServer) sendHeartbeats() {
	heartbeatMsg := &utils.PicoDBServerMessage{
		MessageType:     utils.HEARTBEAT.EnumIndex(),
		SourceIPAddress: server.IPAddress}

	for {
		server.proxyClient.Notify(context.Background(), heartbeatMsg)
		time.Sleep(2 * time.Second)
	}
}

func (server *PicoDBStorageServer) ConstructProxyClient() (
	utils.PicoDBObserverClient, error) {
	var opts []grpc.DialOption
	opts = append(
		opts,
		grpc.WithTransportCredentials(
			insecure.NewCredentials()))

	conn, err := grpc.Dial(server.ProxyServerAddr, opts...)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	serverClient := utils.NewPicoDBObserverClient(conn)

	return serverClient, nil
}

func (server *PicoDBStorageServer) StartServer() {
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

func (server *PicoDBStorageServer) Start() {
	proxyClient, err := server.ConstructProxyClient()

	if err != nil {
		fmt.Println(err)
		return
	}

	server.proxyClient = proxyClient
	go server.sendHeartbeats()

	server.StartServer()
}

func (server *PicoDBStorageServer) Close() {
	server.grpcServer.GracefulStop()
}

func (server *PicoDBStorageServer) ProcessRequest(
	ctx context.Context, request *utils.PicoDBRequest) (*utils.PicoDBResponse, error) {
	response := &utils.PicoDBResponse{}
	response.Key = request.Key

	server.Logger.Log(request)

	switch request.RequestType {
	case utils.GET.EnumIndex():
		value, exists := server.Data[request.Key]

		if !exists {
			response.Value = ""
		} else {
			response.Value = value
		}

	case utils.PUT.EnumIndex():
		server.Data[request.Key] = request.Value
		response.Value = request.Value

	case utils.DELETE.EnumIndex():
		value, exists := server.Data[request.Key]

		if !exists {
			response.Value = ""
		} else {
			response.Value = value
		}

		delete(server.Data, request.Key)
	}

	return response, nil
}
