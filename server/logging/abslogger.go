package logging

import (
	"encoding/json"
	"fmt"
)

type AbsLogger struct {
}

func (logger *AbsLogger) MarshalData(loggerName string, loggerData any) (string, error) {
	bytes, err := json.Marshal(loggerData)

	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s : %s", loggerName, string(bytes)), nil
	}
}
