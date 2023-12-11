package promise

import (
	"net"
	"sync"

	"github.com/SashwatAnagolum/picodb/utils"
	"google.golang.org/protobuf/proto"
)

type PicoDBResultPromise struct {
	IsReady  bool
	Result   *utils.PicoDBResult
	Error    error
	CondVar  *sync.Cond
	ReadLock *sync.Mutex
}

func NewResultPromise() *PicoDBResultPromise {
	return &PicoDBResultPromise{
		IsReady:  false,
		Result:   nil,
		Error:    nil,
		CondVar:  sync.NewCond(&sync.Mutex{}),
		ReadLock: &sync.Mutex{}}
}

func (promise *PicoDBResultPromise) Start(
	serverConn *net.TCPConn, request *utils.PicoDBRequest) {
	promise.CondVar.L.Lock()
	go promise.getResult(serverConn, request)
}

func (promise *PicoDBResultPromise) Ready() bool {
	promise.ReadLock.Lock()
	retVal := promise.IsReady
	promise.ReadLock.Unlock()

	return retVal
}

func (promise *PicoDBResultPromise) WaitForResult() (*utils.PicoDBResult, error) {
	for !promise.IsReady {
		promise.CondVar.Wait()
	}

	if promise.Error != nil {
		return nil, promise.Error
	} else {
		return promise.Result, nil
	}
}

func (promise *PicoDBResultPromise) getResult(
	serverConn *net.TCPConn, request *utils.PicoDBRequest) {
	bytes, err := proto.Marshal(request)

	if err != nil {
		promise.Error = err
		promise.setIsReady()
		return
	}

	_, err = serverConn.Write(bytes)

	if err != nil {
		promise.Error = err
		promise.setIsReady()
		return
	}

	buffer := make([]byte, 2048)
	responseNumBytes, err := serverConn.Read(buffer)

	if err != nil {
		promise.Error = err
		promise.setIsReady()
		return
	} else {
		result := &utils.PicoDBResult{}
		err := proto.Unmarshal(buffer[:responseNumBytes], result)

		if err != nil {
			promise.Error = err
			promise.setIsReady()
			return
		}

		promise.Result = result
		promise.setIsReady()
	}
}

func (promise *PicoDBResultPromise) setIsReady() {
	promise.ReadLock.Lock()
	promise.IsReady = true
	promise.CondVar.Broadcast()
	promise.ReadLock.Unlock()
}
