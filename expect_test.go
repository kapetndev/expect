package expect_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/crumbandbase/expect"
)

const (
	wantError   = true
	wantNoError = false
)

var (
	enterprise = starship{"enterprise"}
	voyager    = starship{"voyager"}
)

type starship struct {
	Name string `json:"name"`
}

type errorDecoder struct{}

func (d *errorDecoder) Decode(interface{}) error {
	return errors.New("bad thing happened")
}

type expectFunc func(*testing.T, interface{}, interface{}, string, ...cmp.Option)
type expectStreamFunc func(*testing.T, expect.Decoder, interface{}, string, ...cmp.Option)

func testFunc(t *testing.T, fn expectFunc, got, expected interface{}, msg string, wantError bool) {
	test := &testing.T{}

	captureOutput(func() {
		fn(test, got, expected, "") // Output is generated in this function.
	})

	// When the test fails but no error was expected, or when the test does not
	// fail but an error was expected.
	if test.Failed() != wantError {
		t.Error(msg)
	}
}

func TestEqual(t *testing.T) {
	t.Run("succeeds when the expected value equals the actual value", func(t *testing.T) {
		testFunc(t, expect.Equal, enterprise, enterprise, "values were not equal", wantNoError)
	})

	t.Run("fails when the expected value does not equal the actual value", func(t *testing.T) {
		testFunc(t, expect.Equal, enterprise, voyager, "values were equal", wantError)
	})
}

func TestNotEqual(t *testing.T) {
	t.Run("succeeds when the expected value does not equal the actual value", func(t *testing.T) {
		testFunc(t, expect.NotEqual, enterprise, voyager, "values were equal", wantNoError)
	})

	t.Run("fails when the expected value equals the actual value", func(t *testing.T) {
		testFunc(t, expect.NotEqual, enterprise, enterprise, "values were not equal", wantError)
	})
}

func newStreamDecoder(t *testing.T, v interface{}) expect.Decoder {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	w := httptest.NewRecorder()
	w.Write(b)

	return json.NewDecoder(w.Body)
}

func testStreamFunc(t *testing.T, fn expectStreamFunc, d expect.Decoder, expected interface{}, msg string, wantError bool) {
	test := &testing.T{}

	captureOutput(func() {
		fn(test, d, expected, "") // Output is generated in this function.
	})

	// When the test fails but no error was expected, or when the test does not
	// fail but an error was expected.
	if test.Failed() != wantError {
		t.Error(msg)
	}
}

func TestStreamEqual(t *testing.T) {
	t.Run("succeeds when the expected response equals the actual response", func(t *testing.T) {
		decoder := newStreamDecoder(t, enterprise)
		testStreamFunc(t, expect.StreamEqual, decoder, &enterprise, "values were not equal", wantNoError)
	})

	t.Run("fails when the expected response does not equal the actual response", func(t *testing.T) {
		decoder := newStreamDecoder(t, enterprise)
		testStreamFunc(t, expect.StreamEqual, decoder, &voyager, "values were equal", wantError)
	})

	t.Run("fails when the expected response is nil", func(t *testing.T) {
		shouldPanic(t, func() {
			decoder := newStreamDecoder(t, enterprise)
			testStreamFunc(t, expect.StreamEqual, decoder, nil, "values were not equal", wantError)
		})
	})

	t.Run("fails when the decoder is nil", func(t *testing.T) {
		shouldPanic(t, func() {
			testStreamFunc(t, expect.StreamEqual, nil, &enterprise, "values were not equal", wantError)
		})
	})

	t.Run("fails when the decoder cannot decode the stream", func(t *testing.T) {
		shouldPanic(t, func() {
			testStreamFunc(t, expect.StreamEqual, &errorDecoder{}, &enterprise, "values were not equal", wantError)
		})
	})
}

func TestStreamNotEqual(t *testing.T) {
	t.Run("succeeds when the expected response does not equal the actual response", func(t *testing.T) {
		decoder := newStreamDecoder(t, enterprise)
		testStreamFunc(t, expect.StreamNotEqual, decoder, &voyager, "values were equal", wantNoError)
	})

	t.Run("fails when the expected response equals the actual response", func(t *testing.T) {
		decoder := newStreamDecoder(t, enterprise)
		testStreamFunc(t, expect.StreamNotEqual, decoder, &enterprise, "values were not equal", wantError)
	})

	t.Run("fails when the expected response is nil", func(t *testing.T) {
		shouldPanic(t, func() {
			decoder := newStreamDecoder(t, enterprise)
			testStreamFunc(t, expect.StreamNotEqual, decoder, nil, "values were not equal", wantError)
		})
	})

	t.Run("fails when the decoder is nil", func(t *testing.T) {
		shouldPanic(t, func() {
			testStreamFunc(t, expect.StreamNotEqual, nil, &enterprise, "values were not equal", wantError)
		})
	})

	t.Run("fails when the decoder cannot decode the stream", func(t *testing.T) {
		shouldPanic(t, func() {
			testStreamFunc(t, expect.StreamNotEqual, &errorDecoder{}, &enterprise, "values were not equal", wantError)
		})
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

func shouldPanic(t *testing.T, fn func()) {
	ch := make(chan bool)
	go func() {
		defer func() { recover(); ch <- true }()
		fn()
		t.Error("did not panic")
	}()
	<-ch
}
