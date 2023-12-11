package utils

type PicoDBServerMessageType int32

const (
	HEARTBEAT PicoDBServerMessageType = iota
	RESPONSE
)

func (reqType PicoDBServerMessageType) String() string {
	return []string{"HEARTBEAT", "RESPONSE"}[reqType]
}

func (reqType PicoDBServerMessageType) EnumIndex() int32 {
	return int32(reqType)
}
