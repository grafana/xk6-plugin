package deno

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_registerPlugin(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		registerPlugin("util.plugin.ts")
	})
}
