package container // import "github.com/docker/docker/container"

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

// Bool marshals to true/false booleans
type Bool struct {
	v int64
}

func (b *Bool) Set(v bool) {
	if v {
		atomic.StoreInt64(&b.v, 1)
	} else {
		atomic.StoreInt64(&b.v, 0)
	}
}

func (b *Bool) Get() bool {
	return atomic.LoadInt64(&b.v) == 1
}

func (b *Bool) UnmarshalJSON(buf []byte) error {
	txt := string(buf)
	switch txt {
	case "true":
		*b = Bool{v: 1}
	case "false":
		*b = Bool{v: 0}
	default:
		return fmt.Errorf("unrecognized: %s", txt)
	}
	return nil
}

func (b Bool) MarshalJSON() ([]byte, error) {
	switch atomic.LoadInt64(&b.v) {
	case 0:
		return []byte("false"), nil
	case 1:
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

type String struct {
	v atomic.Value
}

func (b *String) Set(v string) {
	b.v.Store(v)
}

func (b *String) Get() string {
	s := b.v.Load()
	if s == nil {
		return ""
	}
	return s.(string)
}

func (b *String) UnmarshalJSON(buf []byte) error {
	txt := string(buf)
	var v atomic.Value
	v.Store(txt)
	*b = String{v: v}
	return nil
}

func (b String) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.v.Load())
}
