package encoding

import (
	"encoding/xml"
	"fmt"
	"runtime/debug"
	"testing"
)

type testCodec struct{}

func (c testCodec) Marshal(v interface{}) ([]byte, error) {
	panic("implement me")
}

func (c testCodec) Unmarshal(data []byte, v interface{}) error {
	panic("implement me")
}

func (c testCodec) Name() string {
	return ""
}

// testCodec2 is a Codec implementation with xml.
type testCodec2 struct{}

func (testCodec2) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (testCodec2) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func (testCodec2) Name() string {
	return "xml"
}

func TestRegisterCodec(t *testing.T) {
	f := func() { RegisterCodec(nil) }
	funcDidPanic, panicValue, _ := didPanic(f)
	if !funcDidPanic {
		t.Fatalf(fmt.Sprintf("func should panic\n\tPanic value:\t%#v", panicValue))
	}
	if panicValue != "cannot register a nil Codec" {
		t.Fatalf("panic error got %s want cannot register a nil Codec", panicValue)
	}
	f = func() {
		RegisterCodec(testCodec{})
	}
	funcDidPanic, panicValue, _ = didPanic(f)
	if !funcDidPanic {
		t.Fatalf(fmt.Sprintf("func should panic\n\tPanic value:\t%#v", panicValue))
	}
	if panicValue != "cannot register Codec with empty string result for Name()" {
		t.Fatalf("panic error got %s want cannot register Codec with empty string result for Name()", panicValue)
	}
	codec := testCodec2{}
	RegisterCodec(codec)
	got := GetCodec("xml")
	if got != codec {
		t.Fatalf("RegisterCodec(%v) want %v got %v", codec, codec, got)
	}
}

// PanicTestFunc defines a func that should be passed to the assert.Panics and assert.NotPanics
// methods, and represents a simple func that takes no arguments, and returns nothing.
type PanicTestFunc func()

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f PanicTestFunc) (bool, interface{}, string) {
	didPanic := false
	var message interface{}
	var stack string
	func() {
		defer func() {
			if message = recover(); message != nil {
				didPanic = true
				stack = string(debug.Stack())
			}
		}()

		// call the target function
		f()
	}()

	return didPanic, message, stack
}