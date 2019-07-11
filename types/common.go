package types

import (
	"encoding/hex"
	"reflect"

	"github.com/hosting-de-labs/go-netbox/netbox/models"

	"golang.org/x/crypto/sha3"
)

//HashableEntity has a method to return a string that stays the same when the entity wasn't changed
type HashableEntity interface {
	GetHashableString() string
}

//Metadata contain information that are relevant to communicate with Netbox
type Metadata struct {
	ID           int64
	NetboxEntity interface{}
	EntityType   reflect.Type
}

type NetboxEntity struct {
	entity     interface{}
	entityType reflect.Type
}

func (n NetboxEntity) DcimDevice() models.Device {
	return n.entity.(models.Device)
}

func (n NetboxEntity) DcimDeviceWithConfigContext() models.DeviceWithConfigContext {
	return n.entity.(models.DeviceWithConfigContext)
}

//GetEntity returns the entity and
func (m Metadata) GetEntity() NetboxEntity {
	return NetboxEntity{
		entity:     m.NetboxEntity,
		entityType: m.EntityType,
	}
}

//SetEntity stores a netbox entity with its type
func (m *Metadata) SetEntity(entity interface{}) {
	m.NetboxEntity = entity
	m.EntityType = reflect.TypeOf(entity)
}

//CommonEntity is a general object that should be extended by every Entity that interfaces with Netbox
type CommonEntity struct {
	OriginalEntity interface{}
	Metadata       Metadata
}

//GetIdentifier returns a hash value made from the hashable string
func GetIdentifier(i interface{}) string {
	hashableString := i.(HashableEntity).GetHashableString()

	hash := sha3.Sum512([]byte(hashableString))
	return hex.EncodeToString(hash[:8])
}
