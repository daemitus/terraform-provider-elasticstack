package util

import (
	"bytes"
	"encoding/json"
)

// JsonMarshal is a wrapper for json.Marshal.
func JsonMarshal(val any) ([]byte, error) {
	return json.Marshal(val)
}

// JsonMarshalS is a wrapper for json.Marshal.
func JsonMarshalS(val any) (string, error) {
	v, err := JsonMarshal(val)
	return string(v), err
}

// JsonUnmarshal is a wrapper for json.Unmarshal.
func JsonUnmarshal[T any](val []byte) (T, error) {
	var dest T
	dec := json.NewDecoder(bytes.NewReader(val))
	dec.DisallowUnknownFields()
	err := dec.Decode(&dest)
	return dest, err
}

// JsonUnmarshalS is a wrapper for json.Unmarshal.
func JsonUnmarshalS[T any](val string) (T, error) {
	return JsonUnmarshal[T]([]byte(val))
}
