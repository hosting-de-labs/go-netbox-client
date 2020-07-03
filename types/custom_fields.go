package types

import (
	"fmt"
	"reflect"
)

type CustomFields struct {
	ids    map[string]int
	fields map[string]*string
}

//Load accepts and parses an interface representing go-netbox compatible custom fields
func (c *CustomFields) Load(cf interface{}) (err error) {
	if cf == nil {
		return fmt.Errorf("customfields is nil")
	}

	c.ids = make(map[string]int)
	c.fields = make(map[string]*string)

	for k, f := range cf.(map[string]interface{}) {
		fmt.Printf("%s: %+v\n", k, f)

		switch f.(type) {
		case map[string]interface{}:
			f := f.(map[string]interface{})
			if _, ok := f["value"]; !ok {
				return fmt.Errorf("invalid custom fields: no Value field")
			}

			if _, ok := f["label"]; !ok {
				return fmt.Errorf("invalid custom fields: no Label field")
			}

			c.ids[k] = f["value"].(int)

			val := f["label"].(string)
			c.fields[k] = &val

		case string:
			val := f.(string)
			c.fields[k] = &val

		case nil:
			c.fields[k] = nil

		default:
			return fmt.Errorf("invalid custom fields: field of type %s found", reflect.TypeOf(f))
		}
	}

	return nil
}

func (c *CustomFields) Val(key string) (val *string) {
	f, ok := c.fields[key]
	if !ok {
		return nil
	}

	return f
}

func (c *CustomFields) ValMap() map[string]interface{} {
	if c.fields == nil {
		return nil
	}

	out := make(map[string]interface{}, len(c.fields))

	for k, v := range c.fields {
		if id, ok := c.ids[k]; ok {
			out[k] = id
		} else {
			out[k] = v
		}
	}

	return out
}
