package y

import "reflect"

// Debug enables additional info
var Debug = true

// New creates a proxy of an interface
func New(v interface{}) *Proxy {
	return proxyOf(reflect.ValueOf(v))
}
