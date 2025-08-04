package main

import (
	"fmt"
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


func main(){
	
	args := os.Args

	if(len(args) <= 1){
		println("No arguments")
		return
	}

	switch args[1] {
	case "init":
		createGitDir()
	
	}

}