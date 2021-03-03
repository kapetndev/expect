package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crumbandbase/expect"
)

func TestGreeting(t *testing.T) {
	t.Run("succeeds when the greetings are equal", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := &http.Request{}

		// Call the handler.
		handler(w, r)

		expect.StreamEqual(t, json.NewDecoder(w.Body), &greeting{
			Greeting: "Hello, Picard",
		})
	})

	t.Run("fails when the greeting are not equal", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := &http.Request{}

		// Call the handler.
		handler(w, r)

		expect.StreamEqual(t, json.NewDecoder(w.Body), &greeting{
			Greeting: "Bonjour, Picard",
		})
	})
}
