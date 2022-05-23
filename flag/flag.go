package flag

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
	"google.golang.org/protobuf/proto"

	"github.com/sraphs/go/x/flat"

	"github.com/sraphs/encoding/json"
)

// Name is the name registered for the flag codec.
const Name = "flag"

// Codec is a Codec implementation with flag.
type Codec struct{}

func (Codec) Marshal(v interface{}) ([]byte, error) {
	b, err := json.Codec{}.Marshal(v)
	if err != nil {
		return nil, err
	}

	m := gjson.ParseBytes(b).Value().(map[string]interface{})

	fo := flat.Option{
		Separator: ".",
	}

	f := fo.Flatten(m)

	// sort the keys
	keys := make([]string, 0, len(f))
	for key := range f {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for i, k := range keys {
		if i == len(keys)-1 {
			fmt.Fprintf(&buf, "--%s=%v", k, f[k])
		} else {
			fmt.Fprintf(&buf, "--%s=%v ", k, f[k])
		}
	}

	return buf.Bytes(), nil
}

func (Codec) Unmarshal(data []byte, v interface{}) error {
	s := string(data)

	args := strings.Split(s, " ")

	m, err := Parse(args, v)

	if err != nil {
		return err
	}

	if pm, ok := v.(proto.Message); ok {
		return DecodeValues(pm, m)
	} else if pm, ok := reflect.Indirect(reflect.ValueOf(v)).Interface().(proto.Message); ok {
		return DecodeValues(pm, m)
	}

	fo := flat.Option{
		Separator: ".",
	}

	mi := mapStringToInterface(m)

	uf := fo.Unflatten(mi)

	decoder, err := mapstructure.NewDecoder(defaultDecoderConfig(v))
	if err != nil {
		return err
	}
	err = decoder.Decode(uf)

	return err
}

func (Codec) Name() string {
	return Name
}

// mapStringToInterface converts a map[string]string to a map[string]interface{}
func mapStringToInterface(m map[string]string) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range m {
		out[k] = v
	}
	return out
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
