package encoding

import (
	"errors"
	"strings"

	"github.com/sraphs/go/encoding/env"
	"github.com/sraphs/go/encoding/flag"
	"github.com/sraphs/go/encoding/form"
	"github.com/sraphs/go/encoding/json"
	"github.com/sraphs/go/encoding/proto"
	"github.com/sraphs/go/encoding/xml"
	"github.com/sraphs/go/encoding/yaml"
)

var (
	ErrUnsupportedCodec = errors.New("unsupported codec")
)

// Codec defines the interface Transport uses to encode and decode messages.  Note
// that implementations of this interface must be thread safe; a Codec's
// methods can be called from concurrent goroutines.
type Codec interface {
	// Marshal returns the wire format of v.
	Marshal(v interface{}) ([]byte, error)
	// Unmarshal parses the wire format into v.
	Unmarshal(data []byte, v interface{}) error
	// Name returns the name of the Codec implementation. The returned string
	// will be used as part of content type in transmission.  The result must be
	// static; the result cannot change between calls.
	Name() string
}

var registeredCodecs = make(map[string]Codec)

// RegisterCodec registers the provided Codec for use with all Transport clients and
// servers.
func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	if codec.Name() == "" {
		panic("cannot register Codec with empty string result for Name()")
	}
	contentSubtype := strings.ToLower(codec.Name())
	registeredCodecs[contentSubtype] = codec
}

// GetCodec gets a registered Codec by content-subtype, or nil if no Codec is
// registered for the content-subtype.
//
// The content-subtype is expected to be lowercase.
func GetCodec(contentSubtype string) (codec Codec) {
	if codec, ok := registeredCodecs[contentSubtype]; ok {
		return codec
	}

	switch contentSubtype {
	case "yaml", "yml":
		codec = yaml.Codec{}
	case "xml":
		codec = xml.Codec{}
	case "proto":
		codec = proto.Codec{}
	case "json":
		codec = json.Codec{}
	case "form":
		codec = form.Codec{}
	case "flag":
		codec = flag.Codec{}
	case "env":
		codec = env.Codec{}
	default:
		return unsupportedCodec{}
	}

	RegisterCodec(codec)
	return codec
}

var _ Codec = (*unsupportedCodec)(nil)

type unsupportedCodec struct{}

func (unsupportedCodec) Marshal(v interface{}) ([]byte, error) {
	return nil, ErrUnsupportedCodec
}

func (unsupportedCodec) Unmarshal(data []byte, v interface{}) error {
	return ErrUnsupportedCodec
}

func (unsupportedCodec) Name() string {
	return ""
}
