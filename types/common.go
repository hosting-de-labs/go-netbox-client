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
	entity interface{}
	Meta   *Metadata
}

func (c CommonEntity) GetEntity() interface{} {
	return c.entity
}

func (c *CommonEntity) SetEntity(entity interface{}) {
	c.entity = entity
}

//Meta contain information that are relevant to communicate with Netbox
type Metadata struct {
	ID             int64
	OriginalEntity interface{}
	NetboxEntity   interface{}
	EntityType     reflect.Type
}

//GetNetboxEntity returns the entity and
func (m Metadata) GetNetboxEntity() (entity interface{}, entityType reflect.Type) {
	return m.NetboxEntity, m.EntityType
}

//SetNetboxEntity stores a netbox entity with its type
func (m *Metadata) SetNetboxEntity(entity interface{}) {
	m.NetboxEntity = entity
	m.EntityType = reflect.TypeOf(entity)
}

//GetIdentifier returns a hash value made from the hashable string
func GetIdentifier(i interface{}) string {
	hashableString := i.(HashableEntity).GetHashableString()

	hash := sha3.Sum512([]byte(hashableString))
	return hex.EncodeToString(hash[:8])
}
