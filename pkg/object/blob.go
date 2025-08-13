package object

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ReadBlobObject reads and prints the content of a git object.
func ReadBlobObject(hash string) ([]byte, error) {
	if len(hash) != 40 {
		return nil, fmt.Errorf("invalid object hash provided")
	}
	folderName := hash[0:2]
	fileName := hash[2:]
	filePath := filepath.Join(".git", "objects", folderName, fileName)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("object not found: %w", err)
	}

	zr, err := zlib.NewReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("error decompressing object: %w", err)
	}
	defer zr.Close()

	decompressed, err := io.ReadAll(zr)
	if err != nil {
		return nil, fmt.Errorf("error reading decompressed content: %w", err)
	}

	parts := bytes.SplitN(decompressed, []byte{0}, 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid object format")
	}

	return parts[1], nil
}
