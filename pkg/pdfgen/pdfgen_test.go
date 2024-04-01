package pdfgen_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gosom/toolkit/pkg/pdfgen"
)

func Test_Generate(t *testing.T) {
	ctx := context.Background()

	buff := bytes.NewBuffer(nil)

	b64 := `data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAABlBMVEX///8AAABVwtN+AAABl0lEQVR42uyYvbHEIAyEl3Hg0CW4FEqD0ijFJRASeNg3kuz7m3P65tBYmfGXwIoVEu64446fjYUSbSkIO7CW9FjyBWwApraU1CdWAP1ccgZE9qktzIJtSP1c8ggASQAyh+YYWMgOYGW+PKjBAUthWWNd7Ry+pP3owOlRKexzXUsK7cLEhgbOiIqxpH5RecYGxJZk2wWB5BZznytifqs4vw/IdVRRc9hV2q5L3OELQGQgq6m5RarT5g5ngG56rrBsjbTr+rQgH4CWSurzQD5KCpriDgEGTWHFyG+P3vEBrAUkqziV5jObLPXJF3DUkEUOQEKTO2ZvgHwF0TaJFW+AnkN6K6kOALm8Vl46Zm5R0zm+VBwvgNmSeNRMacTmau2KM0ArjuS29CnlkLt/yD0+oL9069ahqBG/y/3zwIPLfdI/odkjwRdw9ln6mDULkiefO+CYaInCs8htFuQPeE60SDkHE333DEgjJk57MZ8cHnjR3kSnN+DYdNLngU2bvw3uBgeeHqVD9fQ5yvIB3HHHHf8efwEAAP//E+H/lnTjHs4AAAAASUVORK5CYII=`

	html := `<html><body><h1>test</h1><img src="` + b64 + `"></body></html>`
	err := pdfgen.Generate(ctx, buff, []byte(html))
	require.NoError(t, err)

	require.NotEmpty(t, buff.Bytes())
}
