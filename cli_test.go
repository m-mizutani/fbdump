package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCLI(t *testing.T) {
	require.Error(t, Run([]string{"fbdump", "load"}))
}
