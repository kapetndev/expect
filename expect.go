package expect

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equal asserts that two values are identical. If the values are not identical
// then their differences are printed and the current test is marked to have
// failed.
func Equal(t *testing.T, got, expected interface{}, opts ...cmp.Option) {
	if diff := cmp.Diff(got, expected, opts...); diff != "" {
		fmt.Println(diff)
		t.Fail()
	}
}

// NotEqual asserts that two values are not identical. It the values are
// identical then the test is marked to have failed.
func NotEqual(t *testing.T, got, expected interface{}, opts ...cmp.Option) {
	if diff := cmp.Diff(got, expected, opts...); diff == "" {
		t.Error("values are equal")
	}
}

// Decoder is the interface implemented by types that can decode a data stream
// into a structured value.
type Decoder interface {
	Decode(interface{}) error
}

// StreamEqual does a deep comparison of a data stream against an expected
// value. This function uses reflection to determine the type of the expected
// value in order to know how to properly decode the stream.
func StreamEqual(t *testing.T, d Decoder, expected interface{}, opts ...cmp.Option) {
	got := decodeStream(t, d, expected)
	Equal(t, got, expected, opts...)
}

// StreamNotEqual does a deep comparison of a data stream against an expected
// value. This function uses reflection to determine the type of the expected
// value in order to know how to properly decode the stream.
func StreamNotEqual(t *testing.T, d Decoder, expected interface{}, opts ...cmp.Option) {
	got := decodeStream(t, d, expected)
	NotEqual(t, got, expected, opts...)
}

func decodeStream(t *testing.T, d Decoder, expected interface{}) interface{} {
	dpv := reflect.ValueOf(expected)
	if dpv.Kind() != reflect.Ptr {
		t.Fatal("expected not a pointer")
	}
	if dpv.IsNil() {
		t.Fatal("expected pointer is nil")
	}

	// Use the type of the expected value to create a new zero value of the same
	// type. This will be used to unmarshal the body of the response into.
	expectedType := reflect.TypeOf(expected)
	v := reflect.New(expectedType.Elem())

	got := v.Interface()
	if err := d.Decode(got); err != nil {
		t.Fatalf("failed to decode stream: %v", err)
	}

	return got
}
