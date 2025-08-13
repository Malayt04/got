package lstree

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// HandleLsTree reads a tree object and lists its contents.
func HandleLsTree(sha string) {
	folderName := sha[:2]
	fileName := sha[2:]
	filePath := filepath.Join(".git", "objects", folderName, fileName)

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Not a valid object name", sha)
		return
	}

	zr, err := zlib.NewReader(bytes.NewReader(content))
	if err != nil {
		fmt.Println("Error decompressing object:", err)
		return
	}
	defer zr.Close()

	decompressed, err := io.ReadAll(zr)
	if err != nil {
		fmt.Println("Error reading decompressed content:", err)
		return
	}

	parts := bytes.SplitN(decompressed, []byte{0}, 2)
	if len(parts) != 2 {
		fmt.Println("Invalid object format")
		return
	}

	header := string(parts[0])
	body := parts[1]

	if !strings.HasPrefix(header, "tree") {
		fmt.Println("Error: object is not a tree")
		return
	}

	i := 0
	for i < len(body) {
		modeEnd := bytes.IndexByte(body[i:], ' ')
		mode := string(body[i : i+modeEnd])
		i += modeEnd + 1

		nameEnd := bytes.IndexByte(body[i:], 0)
		filename := string(body[i : i+nameEnd])
		i += nameEnd + 1

		shaBytes := body[i : i+20]
		objSha := hex.EncodeToString(shaBytes)
		i += 20

		objType := "blob"
		if mode == "40000" { // Git mode for a tree is "40000" (octal), not "040000"
			objType = "tree"
		}
		fmt.Printf("%s %s %s\t%s\n", mode, objType, objSha, filename)
	}
}