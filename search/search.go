package search

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var cpuCount = runtime.NumCPU()

type fileInfo struct {
	Name string
	Size int64
	Path string
}

type shortFileInfo struct {
	Name string
	Size int64
}

func Scan(path string, removeDuplicates bool) <-chan string {
	_ = removeDuplicates
	fileInfos := walkByDirs(path)

	action := func(path string) {}

	if removeDuplicates {
		action = func(path string) {
			err := os.Remove(path)
			logFatalOnError(err)
		}
	}
	return findDuplicates(fileInfos, action)
}

func walkByDirs(root string) <-chan fileInfo {
	ret := make(chan fileInfo, cpuCount)

	go func() {
		queue := newQueue()
		queue.push(root)

		for queue.len() > 0 {
			dirPath := queue.pop()
			entries, err := os.ReadDir(dirPath)
			logFatalOnError(err)

			for _, entry := range entries {
				if entry.IsDir() {
					queue.push(filepath.Join(dirPath, entry.Name()))
				} else {
					entryInfo, err := entry.Info()
					logFatalOnError(err)

					ret <- fileInfo{
						Name: entry.Name(),
						Path: filepath.Join(dirPath, entry.Name()),
						Size: entryInfo.Size(),
					}
				}
			}
		}
		close(ret)
	}()
	return ret
}

func findDuplicates(fileInfos <-chan fileInfo, actionWithDuplicate func(path string)) <-chan string {
	store := map[shortFileInfo]struct{}{}
	ret := make(chan string, 10)

	go func() {
		for fullFileInfo := range fileInfos {
			info := shortFileInfo{
				Name: fullFileInfo.Name,
				Size: fullFileInfo.Size,
			}
			if _, found := store[info]; found {
				go actionWithDuplicate(fullFileInfo.Path)
				ret <- fullFileInfo.Path
				continue
			}
			store[info] = struct{}{}
		}
		close(ret)
	}()

	return ret
}

func logFatalOnError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}