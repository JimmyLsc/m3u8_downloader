package util

import "github.com/bytedance/sonic"

func GenJsonLog(i interface{}) string {
	result, _ := sonic.MarshalString(i)
	return result
}
