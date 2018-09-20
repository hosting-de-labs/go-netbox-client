package utils

import (
	"reflect"
)

func CompareStruct(item1 interface{}, item2 interface{}, fieldsToCompare []string, fieldsToIgnore []string) bool {
	item1Val := reflect.ValueOf(item1)
	item2Val := reflect.ValueOf(item2)

	item1Type := reflect.TypeOf(item1)
	item2Type := reflect.TypeOf(item2)

	if item1Type.Name() != item2Type.Name() {
		return false
	}

	if item1Type.Kind() != reflect.Struct || item2Type.Kind() != reflect.Struct {
		return false
	}

	if item1Type.Kind() == reflect.Ptr {
		item1Val = item1Val.Elem()
	}

	if item2Type.Kind() == reflect.Ptr {
		item2Val = item2Val.Elem()
	}

LOOP:
	for i := 0; i < item1Val.NumField(); i++ {
		t1 := item1Type.Field(i)

		for _, field := range fieldsToCompare {
			if field != t1.Name {
				continue LOOP
			}
		}

		for _, field := range fieldsToIgnore {
			if field == t1.Name {
				continue LOOP
			}
		}

		if !reflect.DeepEqual(item1Val.Field(i).Interface(), item2Val.Field(i).Interface()) {
			return false
		}
	}

	return true
}
