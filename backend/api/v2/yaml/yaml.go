// Package yaml provides helpers to inspect structs that are tagged with
// `yaml` struct tags.
package yaml

import (
	"fmt"
	"reflect"
)

// GetKeys returns all the YAML keys from a struct that is tagged
// with `yaml` tags.
func GetKeys(s any, parentKey string) []string {
	var tags []string
	if parentKey != "" {
		parentKey = fmt.Sprintf("%s.", parentKey)
	}
	v := reflect.TypeOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := field.Tag.Get("yaml")
		if tag == "" {
			continue
		}
		if field.Type.Kind() == reflect.Struct {
			// This is a new parent.
			newParentKey := fmt.Sprintf("%s%s", parentKey, tag)
			childStruct := reflect.New(field.Type).Interface()
			childTags := GetKeys(childStruct, newParentKey)
			tags = append(tags, childTags...)
		} else {
			// This is a leaf.
			leafTag := fmt.Sprintf("%s%s", parentKey, tag)
			tags = append(tags, leafTag)
		}
	}
	return tags
}
