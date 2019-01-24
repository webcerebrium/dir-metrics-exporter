package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/metrics", indexHandler)

	folder, okFolder := os.LookupEnv("FOLDER")
	if !okFolder {
		panic("Folder for metrics is not specified, please set up FOLDER")
	}
	port, okPort := os.LookupEnv("PORT")
	if !okPort {
		port = "8080"
	}
	host, okH := os.LookupEnv("HOST")
	if !okH {
		host = "0.0.0.0" // default is for docker
	}
	fmt.Println("Listening to " + host + ":" + port + ", watching folder " + folder)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "Please go to /metrics")
}

func getFolderFiles(folder string) int64 {
	return 0
}
func getFolderSize(folder string) int64 {
	return 0
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	folder, _ := os.LookupEnv("FOLDER")
	alias, okAlias := os.LookupEnv("ALIAS")
	if !okAlias {
		alias = "default"
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	io.WriteString(w,
		strings.Join([]string{
			alias+".folder.path "+(folder),
			alias+".folder.files "+string(getFolderFiles(folder)),
			alias+".folder.size "+string(getFolderSize(folder))
		 "\n"))
}
