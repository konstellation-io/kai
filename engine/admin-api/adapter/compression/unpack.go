package compression

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func UnpackFromReader(src io.Reader, dst string) error {
	tmpFile, err := os.CreateTemp("", "process-compress-*.tar.gz")
	if err != nil {
		return fmt.Errorf("creating temp file for process: %w", err)
	}
	defer tmpFile.Close()
	//defer os.Remove(tmpFile.Name())

	_, err = io.Copy(tmpFile, src)
	if err != nil {
		return fmt.Errorf("copying temp file for version: %w", err)
	}

	log.Printf("Created temp file: %s", tmpFile.Name())

	compressedFile, err := os.Open(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("opening process compressed file: %w", err)
	}

	sources, err := gzip.NewReader(compressedFile)
	if err != nil {
		return fmt.Errorf("creating zlib reader: %w", err)
	}
	defer sources.Close()

	tarReader := tar.NewReader(sources)

	for {
		tarFile, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		filePath := filepath.Join(dst, tarFile.Name)

		fmt.Println(filePath)

		if err := processFile(tarReader, filePath, tarFile.Typeflag); err != nil {
			return err
		}
	}

	return nil
}

func processFile(tarReader *tar.Reader, filePath string, fileType byte) error {
	switch fileType {
	case tar.TypeDir:
		if err := os.Mkdir(filePath, 0755); err != nil {
			return fmt.Errorf("error creating krt dir %s: %w", filePath, err)
		}

	case tar.TypeReg:
		outFile, err := os.Create(filePath)

		if err != nil {
			return fmt.Errorf("error creating krt file %s: %w", filePath, err)
		}

		if _, err := io.Copy(outFile, tarReader); err != nil {
			return fmt.Errorf("error copying krt file %s: %w", filePath, err)
		}

		err = outFile.Close()
		if err != nil {
			return fmt.Errorf("error closing krt file %s: %w", filePath, err)
		}

	default:
		return fmt.Errorf("error extracting krt files: uknown type [%v] in [%s]", fileType, filePath)
	}

	return nil
}
