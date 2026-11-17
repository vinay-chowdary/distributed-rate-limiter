package ratelimiter

import (
	"testing"
	"time"
)

func TestAllowConsumesTokens(t *testing.T) {
	// set up
	bucket := NewTokenBucket(3, 1)
	for range 3 {
		result := bucket.Allow()

		expected := true

		if result != expected {
			t.Errorf("got %v, want %v", result, expected)
		}
	}

	result := bucket.Allow()
	expected := false

	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}

}

func TestTokensRefillOverTime(t *testing.T) {
	bucket := NewTokenBucket(2, 1)
	for range 2 {
		result := bucket.Allow()

		expected := true

		if result != expected {
			t.Errorf("got %v, want %v", result, expected)
		}
	}

	result := bucket.Allow()
	expected := false

	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}

	time.Sleep(1 * time.Second)

	result = bucket.Allow()
	expected = true

	if result != expected {
		t.Errorf("got %v, want %v", result, expected)
	}
}
