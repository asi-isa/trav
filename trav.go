package trav

import (
	"io/fs"
	"os"
	"sync"
)

type Trav[T any] struct {
	wg sync.WaitGroup
	ch chan T
}

// Traverse traverses concurrently the file tree rooted at root, calling fn for each file in the tree and onEnd after every file and directory has been traversed.
func (t *Trav[T]) Traverse(root string, fn func(path string, entry fs.DirEntry) (T, bool), onEnd func()) <-chan T {

	clb := func(path string, entry fs.DirEntry) {
		defer t.wg.Done()

		val, ok := fn(path, entry)

		if ok {
			t.ch <- val
		}
	}

	go func() {
		defer onEnd()
		defer t.wg.Wait()

		t.traverse(root, clb)
	}()

	return t.ch
}

func (t *Trav[T]) traverse(path string, fn func(path string, entry fs.DirEntry)) {
	dirEntries, err := os.ReadDir(path)

	if err == nil {
		for _, entry := range dirEntries {
			if entry.IsDir() {
				t.traverse(path+"/"+entry.Name(), fn)
			} else {
				t.wg.Add(1)
				go fn(path, entry)
			}
		}
	}
}
