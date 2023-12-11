package filters

import "github.com/SashwatAnagolum/picodb/utils"

type KVPFilterIF interface {
	GetData() *utils.PicoDBRequest
	SetData(*utils.PicoDBRequest)
}

type AbsKVPFilter struct {
	Data   *utils.PicoDBRequest
	Source KVPFilterIF
}

func (filter *AbsKVPFilter) GetData() *utils.PicoDBRequest {
	if filter.Source != nil {
		return filter.Source.GetData()
	}

	return filter.Data
}

func (filter *AbsKVPFilter) SetData(kvp *utils.PicoDBRequest) {
	if filter.Source == nil {
		filter.Data = kvp
	} else {
		filter.Source.SetData(kvp)
	}
}
