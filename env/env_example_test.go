package env_test

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/tidwall/gjson"

	testData "github.com/sraphs/encoding/internal/testdata/complex"
)

func Example() {
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

	// MarshalOptions is a configurable JSON format marshaller.
	var MarshalOptions = protojson.MarshalOptions{
		EmitUnpopulated: true,
	}

	b, err := MarshalOptions.Marshal(in)

	if err != nil {
		fmt.Println(err)
	}

	m := gjson.ParseBytes(b).Value().(map[string]interface{})

	var UnmarshalOptions = protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	jb, err := json.Marshal(m)

	if err != nil {
		fmt.Println(err)
	}

	var out testData.Complex

	if err := UnmarshalOptions.Unmarshal(jb, &out); err != nil {
		fmt.Println(err)
	}

	fmt.Println(m)
	//Outputs:

}
