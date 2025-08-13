package commit

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/got/pkg/object"
)

// IndexEntry represents a single file in the Git index.
type IndexEntry struct {
	Hash string
	Path string
}

// readIndex reads the simplified text-based index file.
func readIndex() ([]IndexEntry, error) {
	indexPath := ".git/index"
	content, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []IndexEntry{}, nil // Return empty if index doesn't exist
		}
		return nil, err
	}

	var entries []IndexEntry
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			entries = append(entries, IndexEntry{Hash: parts[0], Path: parts[1]})
		}
	}
	return entries, nil
}

// createTreeObject builds a tree object from index entries and writes it.
func createTreeObject(entries []IndexEntry) (string, error) {
	// Sort entries by path for consistent tree objects
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	var treeContent bytes.Buffer
	for _, entry := range entries {
		// For simplicity, we assume all files are blobs with standard permissions.
		// A real implementation handles directories and different file modes.
		mode := "100644" // Regular non-executable file
		path := filepath.Base(entry.Path)
		hashBytes, _ := hex.DecodeString(entry.Hash)

		treeContent.WriteString(fmt.Sprintf("%s %s\x00", mode, path))
		treeContent.Write(hashBytes)
	}

	return object.WriteGitObject("tree", treeContent.Bytes())
}

// getParentCommitHash reads the HEAD to find the current commit hash.
func getParentCommitHash() (string, error) {
	headPath := ".git/HEAD"
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	refPath := strings.TrimSpace(strings.Split(string(headContent), " ")[1])
	fullRefPath := filepath.Join(".git", refPath)

	parentHashBytes, err := os.ReadFile(fullRefPath)
	if os.IsNotExist(err) {
		return "", nil // No parent commit, this is the first commit
	}
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(parentHashBytes)), nil
}

// HandleCommit creates a commit object.
func HandleCommit(message string) {
	// 1. Read the index to get staged files.
	entries, err := readIndex()
	if err != nil {
		fmt.Println("Error reading index:", err)
		return
	}
	if len(entries) == 0 {
		fmt.Println("Nothing to commit, working tree clean.")
		return
	}

	// 2. Create a tree object from the index.
	treeHash, err := createTreeObject(entries)
	if err != nil {
		fmt.Println("Error creating tree object:", err)
		return
	}

	// 3. Get the parent commit hash.
	parentHash, err := getParentCommitHash()
	if err != nil {
		fmt.Println("Error getting parent commit:", err)
		return
	}

	// 4. Construct the commit object content.
	var commitContent strings.Builder
	commitContent.WriteString(fmt.Sprintf("tree %s\n", treeHash))
	if parentHash != "" {
		commitContent.WriteString(fmt.Sprintf("parent %s\n", parentHash))
	}
	// Using placeholder author/committer info
	author := "Go Git User <gogit@example.com>"
	now := time.Now().Format("1136239445 -0700") // Git's internal timestamp format
	commitContent.WriteString(fmt.Sprintf("author %s %s\n", author, now))
	commitContent.WriteString(fmt.Sprintf("committer %s %s\n", author, now))
	commitContent.WriteString(fmt.Sprintf("\n%s\n", message))

	// 5. Write the commit object.
	commitHash, err := object.WriteGitObject("commit", []byte(commitContent.String()))
	if err != nil {
		fmt.Println("Error writing commit object:", err)
		return
	}

	// 6. Update the branch reference (e.g., refs/heads/master) to point to the new commit.
	headContent, _ := os.ReadFile(".git/HEAD")
	refPath := strings.TrimSpace(strings.Split(string(headContent), " ")[1])
	fullRefPath := filepath.Join(".git", refPath)
	if err := os.WriteFile(fullRefPath, []byte(commitHash+"\n"), 0644); err != nil {
		fmt.Println("Error updating branch HEAD:", err)
		return
	}

	// 7. Clear the index after a successful commit.
	if err := os.Remove(".git/index"); err != nil {
		fmt.Println("Warning: could not clear index after commit:", err)
	}

	fmt.Printf("[master (root-commit) %s] %s\n", commitHash[:7], strings.Split(message, "\n")[0])
}