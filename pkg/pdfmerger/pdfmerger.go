package pdfmerger

import (
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/gosom/toolkit/pkg/errorsext"
)

type PDFMerger struct {
	tmpPattern string
}

func New(tmpPattern string) *PDFMerger {
	ans := PDFMerger{}

	if tmpPattern == "" {
		ans.tmpPattern = "pdf_part_*.pdf"
	} else {
		ans.tmpPattern = tmpPattern
	}

	return &ans
}

// Merge merges the input PDFs into a single PDF and writes the result to out.
// The input PDFs are read from in.
// The input PDFs are concatenated in the order they are provided.
// The input PDFs are not modified.
// Merge returns an error if the input PDFs cannot be read or if the output PDF cannot be written.
func (p *PDFMerger) Merge(out io.Writer, in ...io.Reader) error {
	if len(in) == 0 {
		return nil
	}

	if len(in) == 1 {
		if _, err := io.Copy(out, in[0]); err != nil {
			return errorsext.WithStack(err)
		}

		return nil
	}

	tempDir, err := os.MkdirTemp("", "pdfmerge")
	if err != nil {
		return err
	}

	defer os.RemoveAll(tempDir)

	fileNames := make([]string, 0, len(in))

	for _, r := range in {
		tempFile, err := os.CreateTemp(tempDir, p.tmpPattern)
		if err != nil {
			return errorsext.WithStack(err)
		}

		if _, err := io.Copy(tempFile, r); err != nil {
			tempFile.Close()

			return errorsext.WithStack(err)
		}

		fileNames = append(fileNames, tempFile.Name())

		if err := tempFile.Close(); err != nil {
			return errorsext.WithStack(err)
		}
	}

	if err := api.Merge("", fileNames, out, nil, false); err != nil {
		return errorsext.WithStack(err)
	}

	return nil
}
