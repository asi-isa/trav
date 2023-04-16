package trav

import (
	"io/fs"
	"os"
	"sync"
)

type TraverseClbFunc func(path string, entry fs.DirEntry)

func traverse(wg *sync.WaitGroup, path string, fn TraverseClbFunc) {
	dirEntries, err := os.ReadDir(path)

	if err == nil {
		for _, entry := range dirEntries {
			if entry.IsDir() {
				traverse(wg, path+"/"+entry.Name(), fn)
			} else {
				wg.Add(1)
				go fn(path, entry)
			}
		}
	}
}

// Traverse traverses concurrently the file tree rooted at root, calling fn for each file in the tree and onEnd after every file and directory has been traversed.
func Traverse(root string, fn TraverseClbFunc, onEnd func()) {
	var wg sync.WaitGroup

	clb := func(path string, entry fs.DirEntry) {
		defer wg.Done()

		fn(path, entry)
	}

	go func() {
		defer onEnd()
		defer wg.Wait()

		traverse(&wg, root, clb)
	}()
}
