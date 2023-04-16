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
	wg  sync.WaitGroup
	ch  chan DirEntry
	clb func(DirEntry)
}

func Make() Trav {
	return Trav{ch: make(chan DirEntry)}
}

// TODO use predicate to determine, whether the entry should be sent to the channel
func (t *Trav) Traverse(root string, where func(DirEntry) bool) <-chan DirEntry {
	t.clb = func(entry DirEntry) {
		defer t.wg.Done()

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
				go t.clb(makeDirEntry(path, entry))
			}
		}
	}
}
