package cryptoext

import (
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/gosom/toolkit/pkg/errorsext"
)

func GeneratePIN(length int) (string, error) {
	const numLen = 10

	var sb strings.Builder

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(numLen))
		if err != nil {
			return "", errorsext.WithStack(err)
		}

		sb.WriteString(num.String())
	}

	return sb.String(), nil
}
