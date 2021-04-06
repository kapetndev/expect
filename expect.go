package expect

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equal asserts that two values are identical. If the values are not identical
// then their differences are printed and the current test is marked to have
// failed.
func Equal(t *testing.T, got, expected interface{}, msg string, opts ...cmp.Option) {
	if diff := cmp.Diff(got, expected, opts...); diff != "" {
		t.Error(msg)
		t.Logf("\n%s", diff)
	}
}

// NotEqual asserts that two values are not identical. It the values are
// identical then the test is marked to have failed.
func NotEqual(t *testing.T, got, expected interface{}, msg string, opts ...cmp.Option) {
	if diff := cmp.Diff(got, expected, opts...); diff == "" {
		t.Error(msg)
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
func StreamEqual(t *testing.T, d Decoder, expected interface{}, msg string, opts ...cmp.Option) {
	got := decodeStream(t, d, expected)
	Equal(t, got, expected, msg, opts...)
}

// StreamNotEqual does a deep comparison of a data stream against an expected
// value. This function uses reflection to determine the type of the expected
// value in order to know how to properly decode the stream.
func StreamNotEqual(t *testing.T, d Decoder, expected interface{}, msg string, opts ...cmp.Option) {
	got := decodeStream(t, d, expected)
	NotEqual(t, got, expected, msg, opts...)
}

func decodeStream(t *testing.T, d Decoder, expected interface{}) interface{} {
	dpv := reflect.ValueOf(expected)
	if dpv.Kind() != reflect.Ptr || dpv.IsNil() {
		t.Fatal("expected must be a non-nil pointer")
	}

	if d == nil {
		t.Fatal("decoder must be a non-nil type")
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
