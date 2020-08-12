package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	os.Args = []string{"git-remote-cleanup", "--version"}
	version = "0.0.1"
	main()
	assert.Equal(t, 1, 1)
}
