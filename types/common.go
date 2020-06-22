package types

import (
	"fmt"

	"github.com/go-openapi/strfmt"
)

//HashableEntity has a method to return a string that stays the same when the entity wasn't changed
type HashableEntity interface {
	GetHashableString() string
}

//CommonEntity is a general object that should be extended by every Entity that interfaces with Netbox
type CommonEntity struct {
	Meta        *Metadata
	Created     *strfmt.Date
	LastUpdated *strfmt.DateTime
}

//Meta contain information that are relevant to communicate with Netbox
type Metadata struct {
	ID             int64
	CustomFields   CustomFields
	OriginalEntity interface{}
	NetboxEntity   interface{}
}

func (c *CommonEntity) SetCustomFields(customFields interface{}) {
	fmt.Printf("%+v\n\n", customFields)
}

func (c *CommonEntity) SetNetboxEntity(id int64, netboxObj interface{}) {
	if c.Meta == nil {
		c.Meta = &Metadata{}
	}

	c.Meta.ID = id
	c.Meta.NetboxEntity = netboxObj
}

func (c *CommonEntity) SetOriginalEntity(originalObj interface{}) {
	if c.Meta == nil {
		c.Meta = &Metadata{}
	}

	c.Meta.OriginalEntity = originalObj
}

func (c *CommonEntity) HasNetboxEntity() bool {
	return c.Meta != nil && c.Meta.NetboxEntity != nil
}

func (c *CommonEntity) HasOriginalEntity() bool {
	return c.Meta != nil && c.Meta.OriginalEntity != nil
}

func (c *CommonEntity) GetMetaID() int64 {
	if c.Meta == nil {
		return -1
	}

	return c.Meta.ID
}

func (c *CommonEntity) GetMetaOriginalEntity() (out interface{}, ok bool) {
	if c.Meta == nil || c.Meta.OriginalEntity == nil {
		return nil, false
	}

	return c.Meta.OriginalEntity, true
}

func (c *CommonEntity) GetMetaNetboxEntity() (out interface{}, ok bool) {
	if c.Meta == nil || c.Meta.NetboxEntity == nil {
		return nil, false
	}

	return c.Meta.NetboxEntity, true
}
