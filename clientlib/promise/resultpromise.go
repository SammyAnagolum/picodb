package promise

import (
	"context"
	"sync"

	"github.com/SashwatAnagolum/picodb/utils"
)

type PicoDBResponsePromise struct {
	IsReady  bool
	Result   *utils.PicoDBResponse
	Error    error
	CondVar  *sync.Cond
	ReadLock *sync.Mutex
}

func NewResultPromise() *PicoDBResponsePromise {
	return &PicoDBResponsePromise{
		IsReady:  false,
		Result:   nil,
		Error:    nil,
		CondVar:  sync.NewCond(&sync.Mutex{}),
		ReadLock: &sync.Mutex{}}
}

func (promise *PicoDBResponsePromise) Start(
	serverConn *utils.PicoDBServerClient,
	request *utils.PicoDBRequest) {
	promise.CondVar.L.Lock()
	go promise.getResult(serverConn, request)
}

func (promise *PicoDBResponsePromise) Ready() bool {
	promise.ReadLock.Lock()
	retVal := promise.IsReady
	promise.ReadLock.Unlock()

	return retVal
}

func (promise *PicoDBResponsePromise) WaitForResult() (
	*utils.PicoDBResponse, error) {
	for !promise.IsReady {
		promise.CondVar.Wait()
	}

	if promise.Error != nil {
		return nil, promise.Error
	} else {
		return promise.Result, nil
	}
}

func (promise *PicoDBResponsePromise) getResult(
	serverConn *utils.PicoDBServerClient,
	request *utils.PicoDBRequest) {
	response, err := (*serverConn).ProcessRequest(
		context.Background(), request)

	if err != nil {
		promise.Error = err
		promise.setIsReady()
		return
	} else {
		promise.Result = response
		promise.setIsReady()
	}
}

func (promise *PicoDBResponsePromise) setIsReady() {
	promise.ReadLock.Lock()
	promise.IsReady = true
	promise.CondVar.Broadcast()
	promise.ReadLock.Unlock()
}
