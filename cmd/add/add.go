package add

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/got/pkg/object"
)

// IndexEntry represents a single file in the Git index.
type IndexEntry struct {
	Hash string
	Path string
}

// addSingleFile creates a blob object from a single file and adds it to the index.
func addSingleFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("error: could not read file %s: %v\n", filePath, err)
		return
	}
	hashString, err := object.WriteGitObject("blob", content)
	if err != nil {
		fmt.Printf("error: could not create blob for %s: %v\n", filePath, err)
		return
	}

	indexPath := ".git/index"
	f, err := os.OpenFile(indexPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error: could not open index: %v\n", err)
		return
	}
	defer f.Close()

	// Use filepath.ToSlash for consistent path separators in the index.
	if _, err = f.WriteString(fmt.Sprintf("%s %s\n", hashString, filepath.ToSlash(filePath))); err != nil {
		fmt.Printf("error: could not write to index for %s: %v\n", filePath, err)
		return
	}
}

// HandleAdd dispatches adding a single file or a directory.
func HandleAdd(path string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("error: pathspec '%s' did not match any files or directories\n", path)
		return
	}

	if info.IsDir() {
		// Walk the directory and add files.
		err := filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// Skip the .git directory itself.
			if d.Name() == ".git" && d.IsDir() {
				return filepath.SkipDir
			}
			// Only add files, not directories.
			if !d.IsDir() {
				addSingleFile(currentPath)
				fmt.Printf("Added '%s'\n", currentPath)
			}
			return nil
		})
		if err != nil {
			fmt.Printf("error: failed to walk directory %s: %v\n", path, err)
		}
	} else {
		// It's a single file.
		addSingleFile(path)
		fmt.Printf("Added '%s' to staging area.\n", path)
	}
}