# encoding

[![CI](https://github.com/sraphs/encoding/actions/workflows/ci.yml/badge.svg)](https://github.com/sraphs/encoding/actions/workflows/ci.yml)

> The Codec interface to unify the serialization/deserialization logic for processing requests.

## Features

- env
- form
- json
- protobuf
- xml
- yaml

## Install

```bash
go get github.com/sraphs/encoding
```

### Interface

You should implement the following Codec interface for your custom codec.

```go
// Codec interface is for serialization and deserialization, notice that these methods must be thread-safe.
type Codec interface {
    Marshal(v interface{}) ([]byte, error)
    Unmarshal(data []byte, v interface{}) error
    Name() string
}
```
## Usage

### Serialization

```go
// You should manually import this package if you use it directly: 
// import _ "github.com/sraphs/encoding/encoding/json"
jsonCodec := encoding.GetCodec("json")
type user struct {
    Name string
    Age string
    state bool
}
u := &user{
    Name:  "sraph",
    Age:   "2",
    state: false,
}
bytes, _ := jsonCodec.Marshal(u)
fmt.Println(string(bytes))
// output {"Name":"sraph","Age":"2"}
```

### Deserialization

```go
// You should manually import this package if you use it directly: 
// import _ "github.com/sraphs/encoding/encoding/json"
jsonCodec := encoding.GetCodec("json")
type user struct {
    Name string
    Age string
    state bool
}
u := &user{}
jsonCodec.Unmarshal([]byte(`{"Name":"kratos","Age":"2"}`), &u)
fmt.Println(*u)
//output &{kratos 2 false}
```

## Example of Codec Implementation

```go
// https://github.com/sraphs/encoding/blob/main/json/json.go
package json

import (
	"encoding/json"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/sraphs/encoding"
)

// Name is the name registered for the json codec.
const Name = "json"

var (
	// MarshalOptions is a configurable JSON format marshaller.
	MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
	}
	// UnmarshalOptions is a configurable JSON format parser.
	UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}
)

func Init() {
	encoding.RegisterCodec(Codec{})
}

// Codec is a Codec implementation with json.
type Codec struct{}

func (Codec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return MarshalOptions.Marshal(m)
	default:
		return json.Marshal(m)
	}
}

func (Codec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return UnmarshalOptions.Unmarshal(data, m)
	default:
		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return UnmarshalOptions.Unmarshal(data, m)
		}
		return json.Unmarshal(data, m)
	}
}

func (Codec) Name() string {
	return Name
}
```

## Contributing

We alway welcome your contributions :clap:

1.  Fork the repository
2.  Create Feat_xxx branch
3.  Commit your code
4.  Create Pull Request


## CHANGELOG
See [Releases](https://github.com/sraphs/encoding/releases)


## License
[MIT © sraph.com](./LICENSE)
