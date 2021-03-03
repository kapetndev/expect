package main

import (
	"testing"

	"github.com/crumbandbase/expect"
)

func TestGreeting(t *testing.T) {
	t.Run("succeeds when the greetings are equal", func(t *testing.T) {
		expected := "Hello, Picard"
		got := greeting("Picard")

		expect.Equal(t, got, expected)
	})

	t.Run("fails when the greeting are not equal", func(t *testing.T) {
		expected := "Bonjour, Picard"
		got := greeting("Picard")

		expect.Equal(t, got, expected)
	})
}
