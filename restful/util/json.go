package util

import (
	"encoding/json"
	"log"
)

func Struct2Json(input interface{}) []byte {
	body, err := json.Marshal(input)
	if err != nil {
		log.Println("Struct2Json:", err)
		return nil
	}
	return body
}
