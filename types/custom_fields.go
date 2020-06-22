package types

import "fmt"

type CustomFields struct{}

func (cf CustomFields) getMap() interface{} {
	var itf interface{}
	itf = cf

	return itf.(interface{})
}

func (cf *CustomFields) Value(key string) (val interface{}) {
	fmt.Printf("%s\n\n%+v\n\n", key, cf)

	return nil
}
