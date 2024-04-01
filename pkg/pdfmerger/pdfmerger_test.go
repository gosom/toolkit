package pdfmerger_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"toolkit/pkg/pdfmerger"
)

func Test_Merge(t *testing.T) {
	filenames := []string{"testdata/1.pdf", "testdata/2.pdf"}

	files := make([]*os.File, 0, len(filenames))
	readers := make([]io.Reader, 0, len(filenames))

	defer func() {
		for _, file := range files {
			file.Close()
		}
	}()

	for _, filename := range filenames {
		file, err := os.Open(filename)
		require.NoError(t, err)

		files = append(files, file)
		readers = append(readers, file)
	}

	outFile, err := os.CreateTemp("", "merged_*.pdf")
	require.NoError(t, err)

	defer func() {
		outFile.Close()

		os.Remove(outFile.Name())
	}()

	merger := pdfmerger.New("")

	err = merger.Merge(outFile, readers...)
	require.NoError(t, err)
}
