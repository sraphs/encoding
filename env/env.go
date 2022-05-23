package env

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/sraphs/go/x/flat"
)

// Name is the name registered for the env codec.
const Name = "env"

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

// Codec is a Codec implementation with json.
type Codec struct{}

func (Codec) Marshal(v interface{}) (b []byte, err error) {
	if v == nil {
		return []byte{}, nil
	}

	vs := make(map[string]interface{})

	if m, ok := v.(proto.Message); ok {
		ev, err := EncodeValues(m)
		if err != nil {
			return nil, err
		}
		vs = mapStringToInterface(ev)
	} else {
		decoder, err := mapstructure.NewDecoder(defaultDecoderConfig(&vs))

		if err != nil {
			return nil, err
		}

		if err := decoder.Decode(v); err != nil {
			return nil, err
		}

		fo := flat.Option{
			Case:      flat.CaseUpper,
			Separator: "_",
		}

		vs = fo.Flatten(vs)
	}

	keys := make([]string, 0, len(vs))
	for key := range vs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for i, key := range keys {
		v := vs[key]
		if i == len(keys)-1 {
			fmt.Fprintf(&buf, "%s=%v", key, v)
		} else {
			fmt.Fprintf(&buf, "%s=%v\n", key, v)
		}
	}

	return buf.Bytes(), nil
}

func (Codec) Unmarshal(data []byte, v interface{}) error {
	if v == nil {
		return nil
	}

	env, err := godotenv.Unmarshal(string(data))

	if err != nil {
		return err
	}

	if m, ok := v.(proto.Message); ok {
		return DecodeValues(m, env)
	} else if m, ok := reflect.Indirect(reflect.ValueOf(v)).Interface().(proto.Message); ok {
		return DecodeValues(m, env)
	}

	env = lowercaseKeys(env)

	fo := flat.Option{
		Separator: "_",
	}

	unflatted := fo.Unflatten(mapStringToInterface(env))

	decoder, err := mapstructure.NewDecoder(defaultDecoderConfig(v))
	if err != nil {
		return err
	}
	err = decoder.Decode(unflatted)
	return err
}

func (Codec) Name() string {
	return Name
}

// defaultDecoderConfig returns default mapsstructure.DecoderConfig with suppot
// of time.Duration values & string slices
func defaultDecoderConfig(output interface{}) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}

	return c
}

// lowercaseKeys converts a map[string]string to a map[string]string with all keys lowercased
func lowercaseKeys(m map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		out[strings.ToLower(k)] = v
	}
	return out
}

// mapStringToInterface converts a map[string]string to a map[string]interface{}
func mapStringToInterface(m map[string]string) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range m {
		out[k] = v
	}
	return out
}
