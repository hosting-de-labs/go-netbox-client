package types

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

type HashableEntity interface {
	GetHashableString() string
}

type CommonEntity struct {
	OriginalEntity interface{}
}

func GetIdentifier(i interface{}) string {
	hashableString := i.(HashableEntity).GetHashableString()

	hash := sha3.Sum512([]byte(hashableString))
	return hex.EncodeToString(hash[:8])
}
