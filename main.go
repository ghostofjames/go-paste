package main

import (
	// "encoding/base64"
	// "fmt"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	// "time"

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

// func get_id() func() int {
// 	id := 0
// 	return func() int {
// 		id++
// 		return id
// 	}
// }

func get_id() string {
	// TODO: write own unique id filename algorithm
	// now := time.Now().Unix()
	// log.Println(now)
	// encoded := base64.URLEncoding.EncodeToString([]byte(fmt.Sprint(now)))
	// log.Println(encoded)
	// return encoded
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
	filename := get_id()

	// Create new file
	dst, err := os.Create(filepath.Join(config.Folder, filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Write file
	if _, err = io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return file name
	io.WriteString(w, filename+"\n")
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//
	files, err := os.ReadDir(filepath.Join(config.Folder))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		io.WriteString(w, fmt.Sprintf("%s\n", f.Name()))
		// io.WriteString(w, f.Name()+"\n")
	}
}

func main() {
	config = Config{
		Host:   getEnv("HOST", "localhost"),
		Port:   getEnv("PORT", "8000"),
		Folder: getEnv("FOLDER", "files"),
	}
	log.Printf("%+v\n", config)

	os.RemoveAll(config.Folder) // Delete files for testing purposes

	// Setup directory for storing files
	err := os.MkdirAll(config.Folder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", readHandler)

	log.Printf("Listing for requests at http://%s:%s/", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.Host), nil))
}
