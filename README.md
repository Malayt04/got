# Got - A Git Implementation in Go

Got is a simplified implementation of Git written in Go. This project demonstrates the core concepts of Git by implementing key commands from scratch, showing how Git works under the hood.

## What is Got?

Got is an educational project that implements a subset of Git's core functionality. It's not meant to be a full replacement for Git, but rather a learning tool to understand how Git works internally by reimplementing key features in Go.

## Implemented Features

Got currently implements the following Git commands:

1. **init** - Initialize a new Git repository
2. **hash-object** - Compute object ID and optionally create a blob from a file
3. **cat-file** - Provide content or type and size information for repository objects
4. **ls-tree** - List the contents of a tree object
5. **add** - Add file contents to the index
6. **commit** - Record changes to the repository
7. **push** - Simulate pushing changes to a remote repository

## Project Structure

```
got/
├── cmd/
│   ├── add/
│   ├── cat-file/
│   ├── commit/
│   ├── hash-object/
│   ├── initialize/
│   ├── ls-tree/
│   └── push/
├── pkg/
│   └── object/
└── main.go
```

Each command is implemented in its own package under `cmd/`, following Go's standard project layout.

## How It Works

### Git Objects

Got implements Git's core object model, which includes:

1. **Blob Objects** - Store file data
2. **Tree Objects** - Store directory structure and file names
3. **Commit Objects** - Store commit information

All objects are stored in `.git/objects/` directory with the same structure as Git, using SHA-1 hashing and zlib compression.

### Commands

#### init
Creates the basic `.git` directory structure:
```
.git/
├── objects/
│   ├── info/
│   └── pack/
├── refs/
│   ├── heads/
│   └── tags/
├── HEAD
├── config
└── description
```

#### hash-object
Computes the SHA-1 hash of a file and optionally stores it as a blob object in the `.git/objects/` directory.

#### cat-file
Reads and decompresses Git objects from the `.git/objects/` directory and outputs their content.

#### ls-tree
Parses tree objects to display the files and directories they contain.

#### add
Stages files by creating blob objects and adding entries to a simplified index file (`.git/index`).

#### commit
Creates commit objects by:
1. Reading the index to get staged files
2. Creating a tree object from the index
3. Getting the parent commit hash
4. Constructing and writing the commit object
5. Updating branch references

#### push
Provides a simulation of the complex network operations involved in pushing to a remote repository.

## Usage

1. Initialize a repository:
   ```
   got init
   ```

2. Add files to the index:
   ```
   got add <file>
   ```

3. Commit changes:
   ```
   got commit -m "Your commit message"
   ```

4. View object contents:
   ```
   got hash-object <file>
   got cat-file -p <object-hash>
   ```

5. List tree contents:
   ```
   got ls-tree --name-only <tree-hash>
   ```

## Learning Value

This project demonstrates several important concepts:

1. **Git Internals** - How Git stores data as objects with SHA-1 hashes
2. **File System Operations** - Working with directories and files in Go
3. **Data Compression** - Using zlib for object compression
4. **Cryptography** - Using SHA-1 for content addressing
5. **Command Line Interface** - Building a CLI application in Go
6. **Package Organization** - Following Go project structure conventions

## Limitations

This is an educational implementation and has several limitations compared to real Git:

1. No branch management beyond a simple HEAD reference
2. No merge functionality
3. Simplified index implementation
4. No networking code for real push operations
5. Limited error handling and edge case management
6. No support for .gitignore files
7. Simplified tree object creation (doesn't handle nested directories properly)

## Building and Running

To build the project:
```bash
go build -o got .
```

To run:
```bash
./got <command> [<args>]
```

## Future Improvements

Potential enhancements for this project:
1. Implement proper branch management
2. Add diff functionality
3. Implement merge capabilities
4. Add support for .gitignore files
5. Improve error handling and validation
6. Add more Git commands (status, log, checkout, etc.)
7. Implement proper tag management

## License

This project is for educational purposes. Feel free to use it to learn about Git internals.
