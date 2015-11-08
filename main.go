package y

import "reflect"

// Debug enables additional info
var Debug = true

func factory(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr:
		return v.Elem()
	case reflect.Struct:
		return v
	}
	panic("y/main: Y supports ptr on struct or struct only.")
}

// New creates a proxy of an interface
func New(v interface{}) *Proxy {
	return makeProxy(
		factory(reflect.ValueOf(v)),
	)
}
