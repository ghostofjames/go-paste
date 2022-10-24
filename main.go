package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var folder = "files"

func get_id() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}

var next_id = get_id()

func uploadHandler(w http.ResponseWriter, req *http.Request) {

	file, handler, err := req.FormFile("file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer file.Close()

	log.Println(file)

	log.Println(handler)

	// Create new file
	filename := fmt.Sprint("file-", next_id())

	dst, err := os.Create(filepath.Join(folder, filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Write file
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return file name
	io.WriteString(w, filename+"\n")
}

func readHandler(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Path[1:]
	path := filepath.Join(folder, filename)

	dat, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(dat)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(filepath.Join(folder))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		io.WriteString(w, f.Name()+"\n")
	}
}

func main() {

	// setup directory for storing files
	os.RemoveAll(folder) // delete files for testing purposes

	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", readHandler)

	log.Println("Listing for requests at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
