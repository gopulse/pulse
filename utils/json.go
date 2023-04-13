package utils

import "encoding/json"

func ToJSON(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func FromJSON(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}
