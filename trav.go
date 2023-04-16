package trav

import (
	"fmt"
	"io/fs"
)

type TraverseClbFunc func(path string, entry fs.DirEntry)

// Traverse traverses the file tree rooted at root, calling fn for each file in the tree.
func Traverse(root string, fn TraverseClbFunc) {
	fmt.Println("Traverse")
}
