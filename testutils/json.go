package testutils

import (
	"encoding/json"
)

func EncodeIndentedJSON(data interface{}) string {
	b, _ := json.MarshalIndent(data, "", "    ")
	return string(b)
}
