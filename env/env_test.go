package env

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	testData "github.com/sraphs/go/encoding/internal/testdata/complex"
)

var testCases = []struct {
	desc    string
	environ []string
	element interface{}
	data    interface{}
}{
	{
		desc:    "no env vars",
		environ: nil,
		data:    nil,
	},
	{
		desc:    "bool value",
		environ: []string{"FOO=true"},
		element: &struct {
			Foo bool
		}{},
		data: &struct {
			Foo bool
		}{
			Foo: true,
		},
	},
	{
		desc:    "equal",
		environ: []string{"FOO=bar"},
		element: &struct {
			Foo string
		}{},
		data: &struct {
			Foo string
		}{
			Foo: "bar",
		},
	},
	{
		desc:    "multiple bool flags without value",
		environ: []string{"BAR=true", "FOO=true"},
		element: &struct {
			Bar bool
			Foo bool
		}{},
		data: &struct {
			Bar bool
			Foo bool
		}{
			Bar: true,
			Foo: true,
		},
	},
	{
		desc:    "map string",
		environ: []string{"FOO_NAME=bar"},
		element: &struct {
			Foo map[string]interface{}
		}{},
		data: &struct {
			Foo map[string]interface{}
		}{
			Foo: map[string]interface{}{
				"name": "bar",
			},
		},
	},
	{
		desc:    "map struct",
		environ: []string{"FOO_NAME_VALUE=bar"},
		element: &struct {
			Foo map[string]struct{ Value string }
		}{},
		data: &struct {
			Foo map[string]struct{ Value string }
		}{
			Foo: map[string]struct{ Value string }{
				"name": {
					Value: "bar",
				},
			},
		},
	},
	{
		desc:    "map struct with sub-struct",
		environ: []string{"FOO_NAME_BAR_VALUE=bar"},
		element: &struct {
			Foo map[string]struct {
				Bar *struct{ Value string }
			}
		}{},
		data: &struct {
			Foo map[string]struct {
				Bar *struct{ Value string }
			}
		}{
			Foo: map[string]struct {
				Bar *struct{ Value string }
			}{
				"name": {
					Bar: &struct {
						Value string
					}{
						Value: "bar",
					},
				},
			},
		},
	},
	{
		desc:    "map struct with sub-map",
		environ: []string{"FOO_NAME1_BAR_NAME2_VALUE=bar"},
		element: &struct {
			Foo map[string]struct {
				Bar map[string]struct{ Value string }
			}
		}{},
		data: &struct {
			Foo map[string]struct {
				Bar map[string]struct{ Value string }
			}
		}{
			Foo: map[string]struct {
				Bar map[string]struct{ Value string }
			}{
				"name1": {
					Bar: map[string]struct{ Value string }{
						"name2": {
							Value: "bar",
						},
					},
				},
			},
		},
	},
	{
		desc:    "slice",
		environ: []string{"FOO=bar,baz"},
		element: &struct {
			Foo []string
		}{},
		data: &struct {
			Foo []string
		}{
			Foo: []string{"bar", "baz"},
		},
	},
	{
		desc:    "slice struct",
		environ: []string{"FOO[0]_NAME=bar", "FOO[1]_NAME=baz"},
		element: &struct {
			Foo []struct {
				Name string
			}
		}{},
		data: &struct {
			Foo []struct {
				Name string
			}
		}{
			Foo: []struct {
				Name string
			}{{Name: "bar"}, {Name: "baz"}},
		},
	},
	{
		desc:    "slice struct with sub-map",
		environ: []string{"FOO[0]_BAR_NAME1_VALUE=bar1", "FOO[1]_BAR_NAME2_VALUE=bar2"},
		element: &struct {
			Foo []struct {
				Bar map[string]struct{ Value string }
			}
		}{},
		data: &struct {
			Foo []struct {
				Bar map[string]struct{ Value string }
			}
		}{
			Foo: []struct {
				Bar map[string]struct{ Value string }
			}{
				{
					Bar: map[string]struct{ Value string }{
						"name1": {
							Value: "bar1",
						},
					},
				},
				{
					Bar: map[string]struct{ Value string }{
						"name2": {
							Value: "bar2",
						},
					},
				},
			},
		},
	},
}

func TestCodecMarshal(t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			c := Codec{}
			got, err := c.Marshal(tt.data)

			require.NoError(t, err)

			data := []byte(strings.Join(tt.environ, "\n"))

			assert.Equal(t, data, got)
		})
	}
}

func TestCodecUnmarshal(t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			c := Codec{}

			data := []byte(strings.Join(tt.environ, "\n"))

			err := c.Unmarshal(data, tt.element)

			require.NoError(t, err)

			assert.Equal(t, tt.data, tt.element)
		})
	}
}

func TestProtoEncodeDecode(t *testing.T) {
	codec := Codec{}

	in := &testData.Complex{
		Id:      2233,
		NoOne:   "2233",
		Simple:  &testData.Simple{Component: "5566"},
		Simples: []string{"3344", "5566"},
		B:       true,
		Sex:     testData.Sex_woman,
		Age:     18,
		A:       19,
		Count:   3,
		Price:   11.23,
		D:       22.22,
		Byte:    []byte("123"),
		Map:     map[string]string{"sraph": "https://sraph.com/"},

		Timestamp: &timestamppb.Timestamp{Seconds: 20, Nanos: 2},
		Duration:  &durationpb.Duration{Seconds: 120, Nanos: 22},
		Field:     &fieldmaskpb.FieldMask{Paths: []string{"1", "2"}},
		Double:    &wrapperspb.DoubleValue{Value: 12.33},
		Float:     &wrapperspb.FloatValue{Value: 12.34},
		Int64:     &wrapperspb.Int64Value{Value: 64},
		Int32:     &wrapperspb.Int32Value{Value: 32},
		Uint64:    &wrapperspb.UInt64Value{Value: 64},
		Uint32:    &wrapperspb.UInt32Value{Value: 32},
		Bool:      &wrapperspb.BoolValue{Value: false},
		String_:   &wrapperspb.StringValue{Value: "sraph"},
		Bytes:     &wrapperspb.BytesValue{Value: []byte("123")},
	}
	content, err := codec.Marshal(in)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual("a=19\nage=18\nb=true\nbool=false\nbyte=MTIz\nbytes=MTIz\ncount=3\nd=22.22\ndouble=12.33\n"+
		"duration=2m0.000000022s\nfield=1,2\nfloat=12.34\nid=2233\nint32=32\nint64=64\n"+
		"map_sraph=https://sraph.com/\nnumberOne=2233\nprice=11.23\nsex=woman\nsimples=3344,5566\n"+
		"string=sraph\ntimestamp=1970-01-01T00:00:20.000000002Z\nuint32=32\nuint64=64\nverySimple_component=5566", string(content)) {
		t.Errorf("rawpath is not equal to %v", string(content))
	}

	in2 := &testData.Complex{}
	err = codec.Unmarshal(content, in2)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual(int64(2233), in2.Id) {
		t.Errorf("expect %v, got %v", int64(2233), in2.Id)
	}
	if !reflect.DeepEqual("2233", in2.NoOne) {
		t.Errorf("expect %v, got %v", "2233", in2.NoOne)
	}
	if reflect.DeepEqual(in2.Simple, nil) {
		t.Errorf("expect %v, got %v", nil, in2.Simple)
	}
	if !reflect.DeepEqual("5566", in2.Simple.Component) {
		t.Errorf("expect %v, got %v", "5566", in2.Simple.Component)
	}
	if reflect.DeepEqual(in2.Simples, nil) {
		t.Errorf("expect %v, got %v", nil, in2.Simples)
	}
	if !reflect.DeepEqual(len(in2.Simples), 2) {
		t.Errorf("expect %v, got %v", 2, len(in2.Simples))
	}
	if !reflect.DeepEqual("3344", in2.Simples[0]) {
		t.Errorf("expect %v, got %v", "3344", in2.Simples[0])
	}
	if !reflect.DeepEqual("5566", in2.Simples[1]) {
		t.Errorf("expect %v, got %v", "5566", in2.Simples[1])
	}

	if !reflect.DeepEqual("https://sraph.com/", in2.Map["sraph"]) {
		t.Errorf("expect %v, got %v", "https://sraph.com/", in2.Map["sraph"])
	}

}
