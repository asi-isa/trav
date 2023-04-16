package trav

import (
	"io/fs"
	"log"
	"os"
	"sync"
	"time"
)

type DirEntry struct {
	Name  string
	IsDir bool
	Size  int64
	Path  string
	Date  time.Time
}

func makeDirEntry(path string, fsDirEntry fs.DirEntry) DirEntry {
	info, err := fsDirEntry.Info()
	if err != nil {
		log.Println("Error while extrating info from file", err)
	}

	return DirEntry{
		Name:  fsDirEntry.Name(),
		IsDir: fsDirEntry.IsDir(),
		Size:  info.Size(),
		Path:  path,
		Date:  info.ModTime(),
	}
}

type Trav struct {
	wg  *sync.WaitGroup
	ch  chan DirEntry
	clb func(string, fs.DirEntry)
}

func (t Trav) Make() Trav {
	return Trav{}
}

func (t *Trav) Traverse(root string, where func(DirEntry) bool) <-chan DirEntry {
	t.ch = make(chan DirEntry)

	t.clb = func(path string, fsentry fs.DirEntry) {
		defer t.wg.Done()

		entry := makeDirEntry(path, fsentry)

		if where(entry) {
			t.ch <- entry
		}
	}

	go func() {
		defer close(t.ch)
		defer t.wg.Wait()

		t.traverse(root)
	}()

	return t.ch
}

func (t *Trav) traverse(path string) {
	dirEntries, err := os.ReadDir(path)

	if err == nil {
		for _, entry := range dirEntries {
			if entry.IsDir() {
				t.traverse(path + "/" + entry.Name())
			} else {
				t.wg.Add(1)
				go t.clb(path, entry)
			}
		}
	}
}
