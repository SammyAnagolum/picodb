package logging

import (
	"encoding/json"
	"fmt"
)

type AbsLogger struct {
	Name string
}

func (logger *AbsLogger) MarshalData(loggerData any) (string, error) {
	bytes, err := json.Marshal(loggerData)

	if err != nil {
		return "", err
	} else {
		return fmt.Sprintf("%s : %s\n", logger.Name, string(bytes)), nil
	}
}
