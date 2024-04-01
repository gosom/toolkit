package cryptoext_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"toolkit/pkg/cryptoext"
)

func Test_GeneratePIN(t *testing.T) {
	t.Parallel()

	pin, err := cryptoext.GeneratePIN(4)
	require.NoError(t, err)

	require.Equal(t, 4, len(pin))

	for i := 0; i < len(pin); i++ {
		require.True(t, pin[i] >= '0' && pin[i] <= '9')
	}
}
