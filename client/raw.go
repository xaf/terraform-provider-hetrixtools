package hetrixtools

import "encoding/json"

func decodeUntypedJSON(body []byte) (any, error) {
	var result any
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
