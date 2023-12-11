package filters

import (
	"github.com/SashwatAnagolum/picodb/utils"
)

type EncryptKVPFilter struct {
	AbsKVPFilter
}

func (filter *EncryptKVPFilter) Encrypt(str string) string {
	stringChars := make([]byte, len(str))

	for i := 0; i < len(str); i++ {
		if str[i] >= 'a' {
			stringChars[i] = 'a' + ((str[i] + 13 - 'a') % 26)
		} else {
			stringChars[i] = 'A' + ((str[i] + 13 - 'A') % 26)
		}
	}

	return string(stringChars)
}

func (filter *EncryptKVPFilter) GetData() *utils.PicoDBRequest {
	data := filter.AbsKVPFilter.GetData()

	data.Key = filter.Encrypt(data.Key)
	data.Value = filter.Encrypt(data.Value)

	return data
}
