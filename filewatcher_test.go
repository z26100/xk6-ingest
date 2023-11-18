package ingest

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"math"
	"os"
	"path"
	"strings"
	"testing"
	"time"
)

const (
	delay       = 1 * time.Second
	local_path  = "test"
	remote_path = "//sim-db-01.ad.smartinmedia.com/SIM-Share/c.arnheiter/FileWatcherTest"
	file_name   = "test.mp4"
)

func TestFileWatcherRemoteRename(t *testing.T) {
	c := make(chan string)
	defer close(c)
	go watch(remote_path, c)
	time.Sleep(delay)
	os.Rename(fmt.Sprintf("%s/cases/%s", remote_path, file_name), fmt.Sprintf("%s/watch/%s", remote_path, file_name))
	log.Println("sleeping")
	time.Sleep(delay)
	os.Rename(fmt.Sprintf("%s/watch/%s", remote_path, file_name), fmt.Sprintf("%s/cases/%s", remote_path, file_name))
	time.Sleep(delay)
}

func TestFileWatcherLocalRename(t *testing.T) {
	c := make(chan string)
	defer close(c)
	go watch("test/watch", c)
	time.Sleep(delay)
	os.Rename(fmt.Sprintf("%s/cases/%s", local_path, file_name), fmt.Sprintf("%s/watch/%s", local_path, file_name))
	time.Sleep(delay)
	os.Rename(fmt.Sprintf("%s/watch/%s", local_path, file_name), fmt.Sprintf("%s/cases/%s", local_path, file_name))
	time.Sleep(delay)
}
func TestFileWatcherLocalCopy(t *testing.T) {
	c := make(chan string)
	defer close(c)
	go watch("test/watch", c)
	time.Sleep(delay)
	_, err := copy(fmt.Sprintf("%s/cases/%s", local_path, file_name), fmt.Sprintf("%s/watch/%s", local_path, file_name))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(delay)
	err = os.Remove(fmt.Sprintf("%s/watch/%s", local_path, file_name))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(delay)
}

func TestFileWatcherRemoteCopy(t *testing.T) {
	c := make(chan string)
	defer close(c)
	go watch(fmt.Sprintf("%s/%s", remote_path, "/watch"), c)
	time.Sleep(delay)
	log.Printf("Copying file %s\n", file_name)
	_, err := copy(fmt.Sprintf("%s/cases/%s", remote_path, file_name), fmt.Sprintf("%s/watch/%s", remote_path, file_name))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(delay)
	log.Printf("Deleting file\n")
	err = os.Remove(fmt.Sprintf("%s/watch/%s", remote_path, file_name))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(delay)
}

var sizeBuffer = make([]int64, 5)

func avg(buf []int64) int64 {
	var r int64
	for _, i := range buf {
		r = r + i
	}
	return int64(math.Round(float64(r) / float64(len(buf))))
}

var i = 0

func checkFileIsReadable(path string) (bool, error) {
	timeout := 1000
	t := time.NewTicker(delay)

	defer t.Stop()
	for {
		select {
		case <-t.C:
			{

				info, err := os.Stat(path)
				if err != nil {
					return false, err
				}
				if i > timeout {
					return false, err
				}
				newSize := info.Size()
				if newSize != avg(sizeBuffer) {
					sizeBuffer[i%5] = newSize
					i++

					return false, nil
				} else {
					return true, nil
				}

			}
		default:
		}
	}
}

func watch(p string, c chan string) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
	defer watcher.Close()
	// Start listening for events.

	log.Printf("Watching at %s\n", p)
	watcher.Add(p)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op == fsnotify.Write {
				go func() {
					file := path.Clean(strings.Split(event.String(), "\"")[1])
					file = strings.ReplaceAll(file, "\\\\", "/")
					isReady, err := checkFileIsReadable(file)
					if err != nil {
						fmt.Println(err)
						return
					}

					if isReady {
						log.Printf("File %s completely written", file)
					}
				}()
			}
			log.Printf("event: %s %s\n", event.Op.String(), event.Name)
			//c <- event.Name
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
func init() {
	log.SetFlags(log.Lmicroseconds)
}
