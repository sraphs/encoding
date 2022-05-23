package flag

import (
	"testing"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/stretchr/testify/assert"

	testData "github.com/sraphs/encoding/internal/testdata/complex"
)

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
		Field:     &fieldmaskpb.FieldMask{Paths: []string{"a.b", "b.c"}},
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

	expected := "--a=19 --age=18 --b=true --bool=false --byte=MTIz --bytes=MTIz --count=3 --d=22.22 --double=12.33 " +
		"--duration=120.000000022s --field=a.b,b.c --float=12.34 --id=2233 --int32=32 --int64=64 " +
		"--map.sraph=https://sraph.com/ --numberOne=2233 --price=11.23 --sex=woman --simples=3344,5566 " +
		"--string=sraph --timestamp=1970-01-01T00:00:20.000000002Z --uint32=32 --uint64=64 --very_simple.component=5566"

	assert.Equal(t, expected, string(content))

	args := "--a=19 --age=18 --b=true --bool=false --byte=MTIz --bytes=MTIz --count=3 --d=22.22 --double=12.33 " +
		"--duration=120.000000022s --field=a.b,b.c --float=12.34 --id=2233 --int32=32 --int64=64 " +
		"--map.sraph=https://sraph.com/ --numberOne=2233 --price=11.23 --sex=woman --simples=3344,5566 " +
		"--string=sraph --timestamp=1970-01-01T00:00:20.000000002Z --uint32=32 --uint64=64 --very_simple.component=5566"

	in2 := &testData.Complex{}
	err = codec.Unmarshal([]byte(args), in2)

	assert.NoError(t, err)

	assert.Equal(t, in, in2)
}
