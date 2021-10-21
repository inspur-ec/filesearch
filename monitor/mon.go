package monitor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

//type MyWatcher *fsnotify.Watcher
func doev(watcher *fsnotify.Watcher, event fsnotify.Event) {
	switch event.Op {
	case fsnotify.Create:
		watcher.Add(event.Name)
		log.Println("create:", event.Name)
	case fsnotify.Rename, fsnotify.Remove:
		log.Println("remove:", event.Name)
		watcher.Remove(event.Name)
	case fsnotify.Write:
		log.Println("write:", event.Name)
	default:
	}
}
func MonitorFiles(watchdir []string) {
	var err error
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event)
				doev(watcher, event)

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()
	for _, p := range watchdir {
		err = watcher.Add(p)
		if err != nil {
			log.Fatal(err)
		}
		err = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			err = watcher.Add(path)
			if err != nil {
				//log.Fatal(err)
				return err
			}
			return nil
		})
		if err != nil {
			fmt.Printf("walk error [%v]\n", err)
		}
	}
	<-done
}
