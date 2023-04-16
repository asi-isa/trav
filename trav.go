package trav

import (
	"io/fs"
	"log"
	"os"
	"sync"
	"time"
)

type FileEntry struct {
	Name string
	Size int64
	Path string
	Date time.Time
}

func makeFileEntry(path string, fsDirEntry fs.DirEntry) FileEntry {
	info, err := fsDirEntry.Info()
	if err != nil {
		log.Println("Error while extrating info from file", err)
	}

	return FileEntry{
		Name: fsDirEntry.Name(),
		Size: info.Size(),
		Path: path,
		Date: info.ModTime(),
	}
}

type Trav struct {
	wg  sync.WaitGroup
	ch  chan FileEntry
	clb func(string, fs.DirEntry)
}

func New() *Trav {
	return &Trav{}
}

func (t *Trav) Traverse(root string, where func(FileEntry) bool) <-chan FileEntry {
	t.ch = make(chan FileEntry)

	t.clb = func(path string, fsentry fs.DirEntry) {
		defer t.wg.Done()

		entry := makeFileEntry(path, fsentry)

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
