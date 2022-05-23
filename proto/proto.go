// Package proto defines the protobuf codec. Importing this package will
// register the codec.
package proto

import (
	"google.golang.org/protobuf/proto"
)

// Name is the name registered for the proto compressor.
const Name = "proto"

// Codec is a Codec implementation with protobuf. It is the default Codec for Transport.
type Codec struct{}

func (Codec) Marshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (Codec) Unmarshal(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}

func (Codec) Name() string {
	return Name
}
