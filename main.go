package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/teris-io/shortid"
)

type Config struct {
	Host   string
	Port   string
	Folder string
}

var config Config

var fs *FileStore

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		log.Println(exists)
		return value
	}
	return fallback
}

func get_filename() string {
	// TODO: write own unique id filename algorithm
	id, _ := shortid.Generate()
	return id
}

func uploadHandler(w http.ResponseWriter, req *http.Request) {
	file, _, err := req.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename := get_filename()

	err = fs.PutFile(filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, fmt.Sprintf("http://%s:%s/%s\n", config.Host, config.Port, filename))
}

func readHandler(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Path[1:]

	file, err := fs.GetFile(filename)
	if err != nil {
		http.Error(w, "404 - not found", http.StatusNotFound)
		return
	}

	w.Write(file)
}

func main() {
	// Load config
	config = Config{
		Host:   getEnv("HOST", "localhost"),
		Port:   getEnv("PORT", "8000"),
		Folder: getEnv("FOLDER", "files"),
	}

	fs = NewFileStore(config.Folder)

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", readHandler)

	log.Printf("Listing for requests at http://%s:%s/", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil))
}
