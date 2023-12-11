package utils

type PicoDBRequestType int32

const (
	GET PicoDBRequestType = iota
	PUT
	DELETE
)

func (reqType PicoDBRequestType) String() string {
	return []string{"GET", "PUT", "DELETE"}[reqType]
}

func (reqType PicoDBRequestType) EnumIndex() int32 {
	return int32(reqType)
}
