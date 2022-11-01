package main

import (
	"bytes"
	"os"
	"testing"
)

func TestPutAndGet(t *testing.T) {
	testfilename := "testfile"
	testfilecontent := "Hello World"

	tmp, _ := os.MkdirTemp("", "test")
	defer os.RemoveAll(tmp)
	fs := NewFileStore(tmp)

	var testfile bytes.Buffer
	testfile.WriteString(testfilecontent)

	// put file
	err := fs.PutFile(testfilename, &testfile)
	if err != nil {
		t.Error(err)
	}

	// get file
	file, err := fs.GetFile(testfilename)
	if err != nil {
		t.Error(err)
	}
	if string(file) != testfilecontent {
		t.Errorf("FileStore.GetFile() = %v, want %v", file, testfilecontent)
	}
}
