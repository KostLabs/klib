package json

import (
	json "github.com/goccy/go-json"
)

func Encode(input any) ([]byte, error) { //goverifier:ignore:any-type
	return json.Marshal(input)
}

func Decode(input []byte, output any) error { //goverifier:ignore:any-type
	return json.Unmarshal(input, output)
}
