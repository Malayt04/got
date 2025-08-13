package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// WriteGitObject compresses and writes data to the .git/objects directory.
// It returns the SHA-1 hash of the object and any error.
func WriteGitObject(objectType string, content []byte) (string, error) {
	header := fmt.Sprintf("%s %d\x00", objectType, len(content))
	fullContent := append([]byte(header), content...)

	hash := sha1.Sum(fullContent)
	hashHex := hex.EncodeToString(hash[:])

	folderName := hashHex[:2]
	fileName := hashHex[2:]
	folderPath := filepath.Join(".git", "objects", folderName)
	objectPath := filepath.Join(folderPath, fileName)

	if _, err := os.Stat(objectPath); err == nil {
		return hashHex, nil // Object already exists
	}

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return "", fmt.Errorf("error creating object directory: %w", err)
	}

	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)
	if _, err := writer.Write(fullContent); err != nil {
		writer.Close()
		return "", fmt.Errorf("error compressing content: %w", err)
	}
	writer.Close()

	if err := os.WriteFile(objectPath, buf.Bytes(), 0444); err != nil {
		return "", fmt.Errorf("error writing object file: %w", err)
	}

	return hashHex, nil
}
