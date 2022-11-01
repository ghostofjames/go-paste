package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

type FileStore struct {
	Folder string
}

func NewFileStore(path string) *FileStore {
	fs := &FileStore{}
	fs.Folder = path

	os.MkdirAll(config.Folder, os.ModePerm)

	return fs
}

func (fs *FileStore) PutFile(filename string, file io.Reader) error {

	dst, err := os.Create(filepath.Join(fs.Folder, filename))
	if err != nil {
		log.Println(err)
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("created file with filename=%s", filename)
	return nil
}

func (fs *FileStore) GetFile(filename string) ([]byte, error) {

	file, err := os.ReadFile(filepath.Join(fs.Folder, filename))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("read file with filename=%s", filename)
	return file, nil
}

func (fs *FileStore) DeleteFile(filename string) error {

	err := os.Remove(filepath.Join(fs.Folder, filename))
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("removed file with filename=%s", filename)
	return nil
}
