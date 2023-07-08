package util

import (
	"encoding/json"
)

func GenJsonLog(i interface{}) string {
	result, _ := json.Marshal(i)
	return string(result)
}
