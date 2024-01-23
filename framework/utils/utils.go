package utils

import (
	"encoding/json"
	"log"
)

func IsJson(s string) bool {
	var obj struct{}
	if err := json.Unmarshal([]byte(s), &obj); err != nil {
		return false
	}

	return true
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
