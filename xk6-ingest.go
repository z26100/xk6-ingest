package file

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"go.k6.io/k6/js/modules"
	"log"
	"os"
)

func init() {
	modules.Register("k6/x/ingest", new(INGEST))
}

type INGEST struct{}

func (*INGEST) MakeDir(path string) {
	os.Mkdir(path, 0755)
}

func (*INGEST) MakeDirAll(path string) {
	os.MkdirAll(path, 0755)
}

func (*INGEST) Copy(source, to string) {
}

func (*INGEST) Rename(oldPath, newPath string) {
	os.Rename(oldPath, newPath)
}

func (*INGEST) Wait(path string) string {
	c := make(chan string)
	defer close(c)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
	defer watcher.Close()
	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op == fsnotify.Write {
					log.Println("modified file:", event.Name)
					c <- event.Name
					return
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	name := <-c
	return name
}
