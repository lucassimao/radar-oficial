package diarios

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func splitPDFAndConvertToMarkdown(pdfContent []byte) (outputDir string, err error) {

	f, err := os.CreateTemp("", "radar-oficial-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}

	defer f.Close()
	defer os.Remove(f.Name())

	_, err = f.Write(pdfContent)
	if err != nil {
		return "", fmt.Errorf("Failed to write PDF content to temp file: %v", err)
	}

	log.Printf("temp Pdf file: %s", f.Name())

	outputDir, err = os.MkdirTemp("", "radar-oficial-*")
	if err != nil {
		log.Printf("❌ Failed to create temp dir: %v", err)
	}

	log.Printf("output dir: %s", outputDir)

	cmd := exec.Command(
		"python",
		"scripts/split_and_convert_pdf.py",
		f.Name(),
		outputDir,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(
			"split/convert failed: %v\noutput:\n%s",
			err, string(output),
		)
	}
	// on success, output contains anything the script printed
	log.Printf("split script output:\n%s", string(output))

	return outputDir, nil
}
