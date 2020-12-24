package fstraverser

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

type FileInformation struct {
	Path string
	Size int64
}

func TraverseStartingFrom(startFolder string, onEach func(FileInformation)) {
	q := stringQueue{}
	q.push(startFolder)

	for !q.isEmpty() {
		curFolder, _ := q.pop()

		f, err := os.Open(curFolder)
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				continue
			}

			log.Panic(err)
		}

		files, err := f.Readdir(0)
		if err != nil {
			log.Panic(err)
		}

		for _, file := range files {
			if file.IsDir() {
				q.push(filepath.Join(curFolder, file.Name()))
			} else {
				onEach(FileInformation{
					Path: filepath.Join(curFolder, file.Name()),
					Size: file.Size(),
				})
			}
		}
	}
}
