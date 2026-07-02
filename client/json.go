package hetrixtools

import "encoding/json"

func marshalWithoutExtra(value any) (map[string]any, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	result := map[string]any{}
	if err := json.Unmarshal(payload, &result); err != nil {
		return nil, err
	}
	delete(result, "Extra")
	return result, nil
}

func marshalMap(value map[string]any) ([]byte, error) {
	return json.Marshal(value)
}
