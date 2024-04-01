package idgen_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"toolkit/pkg/idgen"
)

func Test_Encode_Decode(t *testing.T) {
	intID := int64(10000)
	sID := idgen.Encode(intID)

	require.NotEmpty(t, sID)
	require.Equal(t, 10, len(sID))

	decoded, err := idgen.Decode(sID)
	require.NoError(t, err)

	require.Equal(t, intID, decoded)

	seen := map[string]bool{}

	for i := 0; i < 10000; i++ {
		s := idgen.Encode(int64(i))

		isSeen := seen[s]
		require.False(t, isSeen)

		seen[s] = true

		decoded, err := idgen.Decode(s)

		require.NoError(t, err)
		require.Equal(t, int64(i), decoded)
	}

	require.Equal(t, 10000, len(seen))
}

func Test_Encode_Specific(t *testing.T) {
	intID := int64(16)
	sID := idgen.Encode(intID)

	require.Equal(t, "s18jaGonyx", sID)
}
