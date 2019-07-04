package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsTime(t *testing.T) {
	tests := map[string]bool{
		"1d":   true,
		"1m":   true,
		"1h":   true,
		"1s":   true,
		"1h1m": true,
		"1h1q": false,
		"1dm":  false,
	}

	for str, value := range tests {
		t.Run(fmt.Sprintf("IsTime/%s", str), func(t *testing.T) {
			v := isTime(str)
			assert.Equal(t, value, v)
		})
	}
}

func TestSubTime(t *testing.T) {
	tm := time.Now()

	tests := map[string]time.Time{
		"1d":     tm.Add(-24 * time.Hour),
		"1m":     tm.Add(-1 * time.Minute),
		"1h":     tm.Add(-1 * time.Hour),
		"1h1m1s": tm.Add(-1 * time.Hour).Add(-1 * time.Minute).Add(-1 * time.Second),
	}

	for str, tvalue := range tests {
		t.Run(fmt.Sprintf("SubTime/%s", str), func(t *testing.T) {
			v := subTime(tm, str)
			assert.Equal(t, tvalue, v)
		})
	}
}

func TestAddTime(t *testing.T) {
	tm := time.Now()

	tests := map[string]time.Time{
		"1d":     tm.Add(24 * time.Hour),
		"1m":     tm.Add(1 * time.Minute),
		"1h":     tm.Add(1 * time.Hour),
		"1h1m1s": tm.Add(1 * time.Hour).Add(1 * time.Minute).Add(1 * time.Second),
	}

	for str, tvalue := range tests {
		t.Run(fmt.Sprintf("AddTime/%s", str), func(t *testing.T) {
			v := addTime(tm, str)
			assert.Equal(t, tvalue, v)
		})
	}
}

func TestParseArgs(t *testing.T) {

	timeStr, str := parseArgs([]string{"server1", "-", "down", "10m"}...)
	assert.Equal(t, "server1 - down", str)
	assert.Equal(t, "10m", timeStr)
}
