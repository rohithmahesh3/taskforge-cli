package output

import (
	"encoding/json"
)

func normalizeStructuredOutput(data interface{}) (interface{}, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var normalized interface{}
	if err := json.Unmarshal(payload, &normalized); err != nil {
		return nil, err
	}

	return normalized, nil
}
