package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
)


func createGitDir() {
	root, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	gitPath := filepath.Join(root, ".git")

	dirs := []string{
		filepath.Join(gitPath, "objects", "info"),
		filepath.Join(gitPath, "objects", "pack"),
		filepath.Join(gitPath, "refs", "heads"),
		filepath.Join(gitPath, "refs", "tags"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0777); err != nil {
			fmt.Println("Error creating dir:", dir, err)
		}
	}

	head := []byte("ref: refs/heads/master\n")
	os.WriteFile(filepath.Join(gitPath, "HEAD"), head, 0644)

	os.WriteFile(filepath.Join(gitPath, "config"), []byte{}, 0644)
	os.WriteFile(filepath.Join(gitPath, "description"), []byte("Unnamed repository; edit this file to name the repository.\n"), 0644)

	fmt.Println("Initialized empty Git repository in", gitPath)
}

func readBlogObject(flag string, commitMessage string) {
	switch flag {
	case "-p":
		if len(commitMessage) < 4 {
			fmt.Println("Invalid object hash")
			return
		}

		folderName := commitMessage[0:2]
		fileName := commitMessage[2:]

		filePath := filepath.Join(".git", "objects", folderName, fileName)

		// Check if file exists
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			fmt.Println("Object not found")
			return
		}

		// Read compressed content
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading object:", err)
			return
		}

		// Decompress using zlib
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

		// Split header and body
		parts := bytes.SplitN(decompressed, []byte{0}, 2)
		if len(parts) != 2 {
			fmt.Println("Invalid object format")
			return
		}

		header := string(parts[0])
		body := string(parts[1])

		fmt.Println("Header:", header)
		fmt.Println()
		fmt.Print(body)

	default:
		fmt.Println("Unsupported flag")
	}
}


func main(){
	
	args := os.Args

	if(len(args) <= 1){
		println("No arguments")
		return
	}

	switch args[1] {
	case "init":
		createGitDir()
	
	case "cat-file":
		readBlogObject(args[2], args[3])
	}

}