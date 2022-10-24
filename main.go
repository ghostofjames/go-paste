package main

import (
	// "encoding/base64"
	// "fmt"
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

var config = Config{
	Host:   getEnv("HOST", "localhost"),
	Port:   getEnv("PORT", "8000"),
	Folder: getEnv("FOLDER", "files"),
}

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

	file, _, err := req.FormFile("file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer file.Close()

	// Create new file
	filename := get_id()

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
	filename := req.URL.Path[1:]
	path := filepath.Join(config.Folder, filename)
	dat, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "404 - not found", http.StatusNotFound)
		return
	}

	w.Write(dat)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(filepath.Join(config.Folder))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		io.WriteString(w, f.Name()+"\n")
	}
}

func main() {
	log.Printf("%+v\n", config)

	os.RemoveAll(config.Folder) // delete files for testing purposes

	// setup directory for storing files
	err := os.MkdirAll(config.Folder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/", readHandler)

	log.Println("Listing for requests at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
