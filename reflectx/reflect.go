package reflectx

import (
	"reflect"
	"runtime"
	"sync"
)

// A FieldInfo is metadata for a struct field
type FieldInfo struct {
	Index    []int
	Path     string
	Field    reflect.StructField
	Zero     reflect.Value
	Name     string
	Options  map[string]string
	Embedded bool
	Children []*FieldInfo
	Parent   *FieldInfo
}

// A StructMap is an index of field metadata for a struct
type StructMap struct {
	Tree  *FieldInfo
	Index []*FieldInfo
	Paths map[string]*FieldInfo
	Names map[string]*FieldInfo
}

// GetByPath ...
func (f StructMap) GetByPath(path string) *FieldInfo {
	return f.Paths[path]
}

// Mapper is a general purpose mapper of names to struct fields. A Mapper
// behaves like most marshallers in the standard library, obeying a field tag
// for name mapping but also providing a basic transform function.
type Mapper struct {
	cache      map[reflect.Type]*StructMap
	tagName    string
	tagMapFunc func(string) string
	mapFunc    func(string) string
	mutex      sync.Mutex
}

// NewMapper returns a new mapper using the tagName as its struct field tag.
// If tagName is the empty string, it is ignored
func NewMapper(tagName string) *Mapper {
	return &Mapper{
		cache:   make(map[reflect.Type]*StructMap),
		tagName: tagName,
	}
}

// TraversalsByName returns a slice of int slices which represent the struct traversals
// for each mapped name. Panics if t is not a struct or Indirectable to a struct
// Returns empty int slice for each name not found.
func (m *Mapper) TraversalsByName(t reflect.Type, names []string) [][]int {
	// r := make([][]int, 0, len(names))

	return nil
}

// TraversalsByNameFunc traverses the mapped names and calls fn with the index of
// each name and the struct travesal represented by that name. Panics if it is not
// a struct or Indirectable to a struct.
func (m *Mapper) TraversalsByNameFunc(t reflect.Type, names []string, fn func(int, []int) error) error {
	t = Deref(t)
	// mustBe(t, reflect.Struct)

	return nil
}

// TypeMap returns a mapping of field strings to int slices representing
// the traversal down the struct to reach the field
func (m *Mapper) TypeMap(t reflect.Type) *StructMap {
	m.mutex.Lock()
	mapping, ok := m.cache[t]
	if !ok {
		m.cache[t] = mapping
	}
	m.mutex.Unlock()

	return mapping
}

// Deref is Indirect for reflect.Types
func Deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

type kinder interface {
	Kind() reflect.Kind
}

func mustBe(v kinder, expected reflect.Kind) {
	if k := v.Kind(); k != expected {
		panic(&reflect.ValueError{Method: methodName(), Kind: k})
	}
}

// methodName returns the caller of the function calling methodName
func methodName() string {
	pc, _, _, _ := runtime.Caller(2)
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown method"
	}

	return f.Name()
}
