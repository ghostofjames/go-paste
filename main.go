package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/teris-io/shortid"
)

type Config struct {
	Host   string
	Port   string
	Folder string
}

var config Config

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

	// Get file from request
	file, _, err := req.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get unique filename
	filename := get_filename()

	// Create new file and write file contents
	dst, err := os.Create(filepath.Join(config.Folder, filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return url to access file
	io.WriteString(w, fmt.Sprintf("http://%s:%s/%s\n", config.Host, config.Port, filename))
}

func readHandler(w http.ResponseWriter, req *http.Request) {
	// Read filename from path
	filename := req.URL.Path[1:]

	// Read file
	dat, err := os.ReadFile(filepath.Join(config.Folder, filename))
	if err != nil {
		http.Error(w, "404 - not found", http.StatusNotFound)
		return
	}

	// Return file content
	w.Write(dat)
}

func main() {
	// Load config
	config = Config{
		Host:   getEnv("HOST", "localhost"),
		Port:   getEnv("PORT", "8000"),
		Folder: getEnv("FOLDER", "files"),
	}
	log.Printf("%+v\n", config)

	// Setup directory for storing files
	err := os.MkdirAll(config.Folder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", readHandler)

	log.Printf("Listing for requests at http://%s:%s/", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil))
}
