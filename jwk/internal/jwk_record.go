package internal

import (
	"github.com/square/go-jose"
)

func JwkJsonWebKeyMarshal(v *jose.JsonWebKey) ([]byte, error) {
	return v.MarshalJSON()
}

func JwkJsonWebKeyUnmarshal(data []byte, v *jose.JsonWebKey) error {
	return v.UnmarshalJSON(data)
}
