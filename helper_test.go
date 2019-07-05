package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceNonAlphanum(t *testing.T) {
	tests := map[string]string{
		"hello":        "hello",
		"hello12345":   "hello12345",
		".-hello-.":    ".-hello-.",
		"9134%$!)23ka": "9134_23ka",
		" hello/ ":     "_hello_",
	}

	for str, value := range tests {
		t.Run(fmt.Sprintf("TestReplaceNonAlphanum/%s", str), func(t *testing.T) {
			v := replaceNonAlphanum(str, "_")
			assert.Equal(t, value, v)
		})
	}
}
