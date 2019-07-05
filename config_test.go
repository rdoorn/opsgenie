package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindConfig(t *testing.T) {
	config, err := findConfig()
	assert.Nil(t, err)
	assert.NotEmpty(t, config.TeamID)
}
