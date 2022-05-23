package form

import (
	"net/url"
	"reflect"

	"github.com/go-playground/form/v4"
	"google.golang.org/protobuf/proto"
)

const (
	// Name is form codec name
	Name = "x-www-form-urlencoded"
)

type Codec struct{}

func (c Codec) Marshal(v interface{}) ([]byte, error) {
	encoder := form.NewEncoder()
	encoder.SetTagName("json")

	var vs url.Values
	var err error
	if m, ok := v.(proto.Message); ok {
		vs, err = EncodeValues(m)
		if err != nil {
			return nil, err
		}
	} else {
		vs, err = encoder.Encode(v)
		if err != nil {
			return nil, err
		}
	}
	for k, v := range vs {
		if len(v) == 0 {
			delete(vs, k)
		}
	}
	return []byte(vs.Encode()), nil
}

func (c Codec) Unmarshal(data []byte, v interface{}) error {
	decoder := form.NewDecoder()
	decoder.SetTagName("json")

	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	if m, ok := v.(proto.Message); ok {
		return DecodeValues(m, vs)
	} else if m, ok := reflect.Indirect(reflect.ValueOf(v)).Interface().(proto.Message); ok {
		return DecodeValues(m, vs)
	}

	return decoder.Decode(v, vs)
}

func (Codec) Name() string {
	return Name
}
