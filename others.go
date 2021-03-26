package main

import (
	"encoding/json"
	"fmt"
)

func replacer(key string, value interface{}) (bool, string) {
	str, ok := value.(string)
	// cut the string
	if ok && len(str) > 30 {
		return true, `"` + str[0:35] + `..."`
	}
	return false, ""
}

func parseAnswer(ans []byte) map[string]interface{} {
	var answer map[string]interface{}
	json.Unmarshal(ans, &answer)
	fmt.Println(answer)
	return answer
}
