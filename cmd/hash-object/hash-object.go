package hashobject

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/got/pkg/object"
)

// CreateHashObject is the implementation for the 'hash-object' command.
func CreateHashObject(write bool, filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if write {
		hash, err := object.WriteGitObject("blob", content)
		if err != nil {
			fmt.Println("Error writing object:", err)
			return
		}
		fmt.Println(hash)
	} else {
		// Just print the hash without writing
		header := fmt.Sprintf("blob %d\x00", len(content)) // Corrected: \x00 is a valid escape sequence for a null byte
		fullContent := append([]byte(header), content...)
		hash := sha1.Sum(fullContent)
		fmt.Printf("%x\n", hash)
	}
}
