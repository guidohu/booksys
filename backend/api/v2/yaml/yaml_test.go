package yaml

import (
	"reflect"
	"testing"
)

// The types below only exist to exercise GetKeys. They cover the shapes that
// show up in the configuration structs: leaves, nested structs, dotted tags
// and fields that carry no `yaml` tag at all.

type leaves struct {
	Host    string `yaml:"host"`
	Port    uint   `yaml:"port"`
	Enabled bool   `yaml:"enabled"`
}

type nested struct {
	Leaves leaves `yaml:"leaves"`
	Name   string `yaml:"name"`
}

type deeplyNested struct {
	Nested nested `yaml:"nested"`
}

type partiallyTagged struct {
	Tagged   string `yaml:"tagged"`
	Untagged string
	JSONOnly string `json:"jsonOnly"`
	Skipped  leaves // No yaml tag, so the whole sub tree is skipped.
	Kept     string `yaml:"kept"`
}

type untaggedLeaves struct {
	A string
	B string
}

type withUntaggedSubStruct struct {
	Empty untaggedLeaves `yaml:"empty"`
	Kept  string         `yaml:"kept"`
}

type withPointer struct {
	Child *leaves `yaml:"child"`
	Name  string  `yaml:"name"`
}

type withCollections struct {
	Items []string          `yaml:"items"`
	Meta  map[string]string `yaml:"meta"`
	Fixed [3]int            `yaml:"fixed"`
}

type withDottedTag struct {
	APIKey string `yaml:"api.key"`
}

type withTagOptions struct {
	Name string `yaml:"name,omitempty"`
	Skip string `yaml:"-"`
}

type withEmbeddedTagged struct {
	leaves `yaml:"embedded"`
	Extra  string `yaml:"extra"`
}

type withEmbeddedUntagged struct {
	leaves
	Extra string `yaml:"extra"`
}

type withUnexported struct {
	hidden  string `yaml:"hidden"`
	Visible string `yaml:"visible"`
}

type empty struct{}

// configLike mirrors the layout of config.Configuration, which is the actual
// caller of GetKeys.
type configLike struct {
	Database   leaves        `yaml:"database"`
	MyNautique withDottedTag `yaml:"mynautique"`
}

// equalKeys compares two key lists and treats a nil slice and an empty slice
// as equal.
func equalKeys(a, b []string) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return reflect.DeepEqual(a, b)
}

func TestGetKeys(t *testing.T) {
	tests := []struct {
		name      string
		value     any
		parentKey string
		want      []string
	}{
		{
			name:  "leaves are returned in declaration order",
			value: &leaves{},
			want:  []string{"host", "port", "enabled"},
		},
		{
			name:      "parent key is prefixed",
			value:     &leaves{},
			parentKey: "database",
			want:      []string{"database.host", "database.port", "database.enabled"},
		},
		{
			name:  "nested structs are flattened",
			value: &nested{},
			want:  []string{"leaves.host", "leaves.port", "leaves.enabled", "name"},
		},
		{
			name:      "nested structs keep the parent key",
			value:     &deeplyNested{},
			parentKey: "root",
			want: []string{
				"root.nested.leaves.host",
				"root.nested.leaves.port",
				"root.nested.leaves.enabled",
				"root.nested.name",
			},
		},
		{
			name:  "fields without a yaml tag are skipped",
			value: &partiallyTagged{},
			want:  []string{"tagged", "kept"},
		},
		{
			name:  "a sub struct without yaml tags contributes no key",
			value: &withUntaggedSubStruct{},
			want:  []string{"kept"},
		},
		{
			name:  "a pointer to a struct is treated as a leaf",
			value: &withPointer{},
			want:  []string{"child", "name"},
		},
		{
			name:  "slices, maps and arrays are leaves",
			value: &withCollections{},
			want:  []string{"items", "meta", "fixed"},
		},
		{
			name:  "a dotted tag is kept as is",
			value: &withDottedTag{},
			want:  []string{"api.key"},
		},
		{
			name:      "a dotted tag is prefixed with the parent key",
			value:     &withDottedTag{},
			parentKey: "mynautique",
			want:      []string{"mynautique.api.key"},
		},
		{
			// GetKeys does not parse the tag, so options and the `-` that
			// tells a yaml marshaller to skip a field end up in the key.
			name:  "tag options are not stripped",
			value: &withTagOptions{},
			want:  []string{"name,omitempty", "-"},
		},
		{
			name:  "a tagged embedded struct becomes a parent",
			value: &withEmbeddedTagged{},
			want:  []string{"embedded.host", "embedded.port", "embedded.enabled", "extra"},
		},
		{
			// An embedded struct is inlined by yaml marshallers, but without a
			// tag GetKeys skips it just like any other untagged field.
			name:  "an untagged embedded struct is skipped",
			value: &withEmbeddedUntagged{},
			want:  []string{"extra"},
		},
		{
			name:  "unexported fields are included",
			value: &withUnexported{},
			want:  []string{"hidden", "visible"},
		},
		{
			name:  "an empty struct has no keys",
			value: &empty{},
			want:  nil,
		},
		{
			name:  "an untagged struct has no keys",
			value: &untaggedLeaves{},
			want:  nil,
		},
		{
			// This is how config.GetKeys calls it: only the type matters, the
			// pointer is never dereferenced.
			name:  "a nil typed pointer works",
			value: (*nested)(nil),
			want:  []string{"leaves.host", "leaves.port", "leaves.enabled", "name"},
		},
		{
			name:  "a config like struct",
			value: (*configLike)(nil),
			want: []string{
				"database.host",
				"database.port",
				"database.enabled",
				"mynautique.api.key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetKeys(tt.value, tt.parentKey)
			if !equalKeys(got, tt.want) {
				t.Errorf("GetKeys(%T, %q) = %v, want %v", tt.value, tt.parentKey, got, tt.want)
			}
		})
	}
}

// TestGetKeysPanics documents that GetKeys requires a pointer to a struct.
func TestGetKeysPanics(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{name: "struct value instead of a pointer", value: leaves{}},
		{name: "pointer to a non struct", value: new(string)},
		{name: "untyped nil", value: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("GetKeys(%#v, \"\") did not panic, want a panic", tt.value)
				}
			}()
			_ = GetKeys(tt.value, "")
		})
	}
}
