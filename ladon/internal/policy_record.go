package internal

import (
	"encoding/json"
	"github.com/ory/ladon"
)

func PolicyMarshal(v *ladon.Policy) ([]byte, error) {
	return json.Marshal(v)
}

func PolicyUnmarshal(data []byte, v *ladon.DefaultPolicy) error {
	return v.UnmarshalJSON(data)
}
