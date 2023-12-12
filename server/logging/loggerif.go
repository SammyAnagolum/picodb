package logging

import "github.com/SashwatAnagolum/picodb/utils"

type LoggerIF interface {
	Log(request *utils.PicoDBRequest)
	GetLogs() (string, error)
	ResetLogs()
	GetName() string
}
