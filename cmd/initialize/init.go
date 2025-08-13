package initialize

import (
	"fmt"
	"os"
	"path/filepath"
)

// createGitDir initializes a new .git directory structure.
func CreateGitDir() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	gitPath := filepath.Join(root, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		fmt.Println("Reinitialized existing Git repository in", gitPath)
	} else {
		fmt.Println("Initialized empty Git repository in", gitPath)
	}

	dirs := []string{
		filepath.Join(gitPath, "objects", "info"),
		filepath.Join(gitPath, "objects", "pack"),
		filepath.Join(gitPath, "refs", "heads"),
		filepath.Join(gitPath, "refs", "tags"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("Error creating dir:", dir, err)
		}
	}

	head := []byte("ref: refs/heads/master\n")
	os.WriteFile(filepath.Join(gitPath, "HEAD"), head, 0644)

	config := []byte(`[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
`)
	os.WriteFile(filepath.Join(gitPath, "config"), config, 0644)

	description := []byte("Unnamed repository; edit this file to name the repository.\n")
	os.WriteFile(filepath.Join(gitPath, "description"), description, 0644)
}
