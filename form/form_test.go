package form

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	testData "github.com/sraphs/encoding/internal/testdata/complex"
)

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type TestModel struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

const contentType = "x-www-form-urlencoded"

func TestFormCodecMarshal(t *testing.T) {
	codec := Codec{}

	req := &LoginRequest{
		Username: "sraph",
		Password: "sraph_pwd",
	}
	content, err := codec.Marshal(req)
	if err != nil {
		t.Errorf("marshal error: %v", err)
	}
	if !reflect.DeepEqual([]byte("password=sraph_pwd&username=sraph"), content) {
		t.Errorf("expect %v, got %v", []byte("password=sraph_pwd&username=sraph"), content)
	}

	req = &LoginRequest{
		Username: "sraph",
		Password: "",
	}
	content, err = codec.Marshal(req)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual([]byte("username=sraph"), content) {
		t.Errorf("expect %v, got %v", []byte("username=sraph"), content)
	}

	m := &TestModel{
		ID:   1,
		Name: "sraph",
	}
	content, err = codec.Marshal(m)
	t.Log(string(content))
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual([]byte("id=1&name=sraph"), content) {
		t.Errorf("expect %v, got %v", []byte("id=1&name=sraph"), content)
	}
}

func TestFormCodecUnmarshal(t *testing.T) {
	codec := Codec{}

	req := &LoginRequest{
		Username: "sraph",
		Password: "sraph_pwd",
	}
	content, err := codec.Marshal(req)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}

	bindReq := new(LoginRequest)
	err = codec.Unmarshal(content, bindReq)
	if err != nil {
		t.Errorf("expect %v, got %v", nil, err)
	}
	if !reflect.DeepEqual("sraph", bindReq.Username) {
		t.Errorf("expect %v, got %v", "sraph", bindReq.Username)
	}
	if !reflect.DeepEqual("sraph_pwd", bindReq.Password) {
		t.Errorf("expect %v, got %v", "sraph_pwd", bindReq.Password)
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
	if !reflect.DeepEqual("a=19&age=18&b=true&bool=false&byte=MTIz&bytes=MTIz&count=3&d=22.22&double=12.33&duration="+
		"2m0.000000022s&field=1%2C2&float=12.34&id=2233&int32=32&int64=64&map%5Bsraph%5D=https%3A%2F%2Fsraph.com%2F&"+
		"numberOne=2233&price=11.23&sex=woman&simples=3344&simples=5566&string=sraph"+
		"&timestamp=1970-01-01T00%3A00%3A20.000000002Z&uint32=32&uint64=64&very_simple.component=5566", string(content)) {
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
}
