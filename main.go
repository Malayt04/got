package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"crypto/sha1"
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

		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			fmt.Println("Object not found")
			return
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading object:", err)
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
		body := string(parts[1])

		fmt.Println("Header:", header)
		fmt.Println()
		fmt.Print(body)

	default:
		fmt.Println("Unsupported flag")
	}
}

func createHashObject(flag string, filePath string) {

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}


	header := fmt.Sprintf("blob %d\u0000", len(content))
	store := append([]byte(header), content...)


	hash := sha1.Sum(store)
	hashHex := fmt.Sprintf("%x", hash[:]) 

	fmt.Println(hashHex)

	if flag == "-w" {

		folderName := hashHex[:2]
		fileName := hashHex[2:]

		folderPath := filepath.Join(".git", "objects", folderName)
		if err := os.MkdirAll(folderPath, 0777); err != nil {
			fmt.Println("Error creating folder:", err)
			return
		}

		var buf bytes.Buffer
		writer := zlib.NewWriter(&buf)
		if _, err := writer.Write(store); err != nil {
			fmt.Println("Error compressing content:", err)
			return
		}
		writer.Close()

		objectPath := filepath.Join(folderPath, fileName)
		if err := os.WriteFile(objectPath, buf.Bytes(), 0644); err != nil {
			fmt.Println("Error writing object file:", err)
			return
		}
	}
}

func handleLsTree(flag string, sha string) {
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
		fmt.Println("Not a tree object:", header)
		return
	}

	i := 0
	for i < len(body) {
		// Extract mode
		modeEnd := bytes.IndexByte(body[i:], ' ')
		mode := string(body[i : i+modeEnd])
		i += modeEnd + 1

		// Extract filename
		nameEnd := bytes.IndexByte(body[i:], 0)
		filename := string(body[i : i+nameEnd])
		i += nameEnd + 1

		// Extract SHA (20 bytes)
		if i+20 > len(body) {
			fmt.Println("Unexpected end of object while reading SHA")
			return
		}
		shaBytes := body[i : i+20]
		sha := fmt.Sprintf("%x", shaBytes)
		i += 20

		fmt.Printf("%s %s %s\n", mode, "blob/tree", filename, sha)
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
	
	case "hash-object":
		if(args[2] != "-w"){
			createHashObject("", args[2])
		}
		createHashObject(args[2], args[3])

	case "ls-Tree":
		if(len(args) < 4){
			handleLsTree("", args[2])
		}
		handleLsTree(args[2], args[3])
	}

}