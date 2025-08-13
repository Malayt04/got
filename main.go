package main

import (
	"fmt"
	"os"

	"github.com/got/cmd/add"
	"github.com/got/cmd/cat-file"
	"github.com/got/cmd/commit"
	"github.com/got/cmd/hash-object"
	"github.com/got/cmd/initialize"
	"github.com/got/cmd/ls-tree"
	"github.com/got/cmd/push"
)



func main() {
	// For direct execution, we just need the command as the first argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: got <command> [<args>]")
		return
	}

	// The command is now the first argument (index 1)
	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "init":
		initialize.CreateGitDir()
	case "cat-file":
		if len(args) < 2 || args[0] != "-p" {
			fmt.Println("Usage: got cat-file -p <hash>")
			return
		}
		catfile.ReadBlobObject(args[1])
	case "hash-object":
		if len(args) == 1 {
			hashobject.CreateHashObject(false, args[0])
		} else if len(args) == 2 && args[0] == "-w" {
			hashobject.CreateHashObject(true, args[1])
		} else {
			fmt.Println("Usage: got hash-object [-w] <file>")
		}
	case "ls-tree":
		if len(args) < 2 || args[0] != "--name-only" {
			fmt.Println("Usage: got ls-tree --name-only <hash>")
			return
		}
		lstree.HandleLsTree(args[1])
	case "add":
		if len(args) < 1 {
			fmt.Println("Usage: got add <file-or-directory>")
			return
		}
		add.HandleAdd(args[0])
	case "commit":
		if len(args) < 2 || args[0] != "-m" {
			fmt.Println("Usage: got commit -m \"<message>\"")
			return
		}
		commit.HandleCommit(args[1])
	case "push":
		if len(args) < 2 {
			fmt.Println("Usage: got push <remote> <branch>")
			return
		}
		push.HandlePush(args[0], args[1])
	default:
		fmt.Println("Unknown command:", command)
	}
}
