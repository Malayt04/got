package catfile

import (
	"fmt"

	"github.com/got/pkg/object"
)

// ReadBlobObject reads and prints the content of a git object.
func ReadBlobObject(hash string) {
	content, err := object.ReadBlobObject(hash)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Print(string(content))
}