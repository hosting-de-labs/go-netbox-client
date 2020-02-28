package types

import (
	"encoding/hex"
	"reflect"

	"golang.org/x/crypto/sha3"
)

//HashableEntity has a method to return a string that stays the same when the entity wasn't changed
type HashableEntity interface {
	GetHashableString() string
}

//CommonEntity is a general object that should be extended by every Entity that interfaces with Netbox
type CommonEntity struct {
	Metadata *Metadata
}

//Metadata contain information that are relevant to communicate with Netbox
type Metadata struct {
	ID           int64
	NetboxEntity interface{}
	EntityType   reflect.Type
}

//GetEntity returns the entity and
func (m Metadata) GetEntity() (entity interface{}, entityType reflect.Type) {
	return m.NetboxEntity, m.EntityType
}

//SetEntity stores a netbox entity with its type
func (m *Metadata) SetEntity(entity interface{}) {
	m.NetboxEntity = entity
	m.EntityType = reflect.TypeOf(entity)
}

//GetIdentifier returns a hash value made from the hashable string
func GetIdentifier(i interface{}) string {
	hashableString := i.(HashableEntity).GetHashableString()

	hash := sha3.Sum512([]byte(hashableString))
	return hex.EncodeToString(hash[:8])
}
