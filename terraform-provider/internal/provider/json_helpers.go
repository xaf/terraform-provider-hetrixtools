package provider

import (
	"bytes"
	"encoding/json"
)

func normalizeJSON(body []byte) string {
	var value any
	if err := json.Unmarshal(body, &value); err != nil {
		return string(body)
	}
	formatted, err := json.Marshal(value)
	if err != nil {
		return string(bytes.TrimSpace(body))
	}
	return string(formatted)
}
