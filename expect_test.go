package expect_test

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/crumbandbase/expect"
	"github.com/google/go-cmp/cmp"
)

var (
	enterprise = starship{"enterprise"}
	voyager    = starship{"voyager"}
)

type starship struct {
	Name string `json:"name"`
}

type expectFunc func(*testing.T, interface{}, interface{}, ...cmp.Option)
type expectStreamFunc func(*testing.T, expect.Decoder, interface{}, ...cmp.Option)

func testFunc(fn expectFunc, got, expected interface{}) *testing.T {
	t := &testing.T{}

	captureOutput(func() {
		fn(t, got, expected) // Output is generated in this function.
	})

	return t
}

func TestEqual(t *testing.T) {
	t.Run("succeeds when the expected value equals the actual value", func(t *testing.T) {
		if test := testFunc(expect.Equal, enterprise, enterprise); test.Failed() {
			t.Error("values were not equal")
		}
	})

	t.Run("fails when the expected value does not equal the actual value", func(t *testing.T) {
		if test := testFunc(expect.Equal, enterprise, voyager); !test.Failed() {
			t.Error("values were equal")
		}
	})
}

func TestNotEqual(t *testing.T) {
	t.Run("succeeds when the expected value does not equal the actual value", func(t *testing.T) {
		if test := testFunc(expect.NotEqual, enterprise, voyager); test.Failed() {
			t.Error("values were equal")
		}
	})

	t.Run("fails when the expected value equals the actual value", func(t *testing.T) {
		if test := testFunc(expect.NotEqual, enterprise, enterprise); !test.Failed() {
			t.Error("values were not equal")
		}
	})
}

func testStreamFunc(fn expectStreamFunc, in []byte, expected interface{}) *testing.T {
	t := &testing.T{}
	w := httptest.NewRecorder()
	w.Write(in)

	captureOutput(func() {
		fn(t, json.NewDecoder(w.Body), expected)
	})

	return t
}

func TestStreamEqual(t *testing.T) {
	t.Run("succeeds when the expected response equals the actual response", func(t *testing.T) {
		b, err := json.Marshal(enterprise)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		if test := testStreamFunc(expect.StreamEqual, b, &enterprise); test.Failed() {
			t.Error("values were not equal")
		}
	})

	t.Run("fails when the expected response does not equal the actual response", func(t *testing.T) {
		b, err := json.Marshal(enterprise)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		if test := testStreamFunc(expect.StreamEqual, b, &voyager); !test.Failed() {
			t.Error("values were equal")
		}
	})
}

func TestStreamNotEqual(t *testing.T) {
	t.Run("succeeds when the expected response does not equal the actual response", func(t *testing.T) {
		b, err := json.Marshal(enterprise)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		if test := testStreamFunc(expect.StreamNotEqual, b, &voyager); test.Failed() {
			t.Error("values were equal")
		}
	})

	t.Run("fails when the expected response equals the actual response", func(t *testing.T) {
		b, err := json.Marshal(enterprise)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		if test := testStreamFunc(expect.StreamNotEqual, b, &enterprise); !test.Failed() {
			t.Error("values were not equal")
		}
	})
}

func captureOutput(fn func()) {
	old := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
}
