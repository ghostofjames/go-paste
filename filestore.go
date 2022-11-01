package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
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

func (fs *FileStore) PutFile(name string, file multipart.File) error {
	path := filepath.Join(fs.Folder, name)
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStore) GetFile(name string) ([]byte, error) {
	log.Println(fs.Folder)
	path := filepath.Join(fs.Folder, name)
	log.Println(path)
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("file with name=%s not found", name)
	}

	return file, nil
}

func (fs *FileStore) DeleteFile(name string) error {
	path := filepath.Join(fs.Folder, name)
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
