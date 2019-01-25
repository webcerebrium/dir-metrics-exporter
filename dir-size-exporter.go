package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/", indexHandler)

	folder, okFolder := os.LookupEnv("FOLDER")
	if !okFolder {
		panic("Folder for metrics is not specified, please set up FOLDER")
	}
	port, okPort := os.LookupEnv("PORT")
	if !okPort {
		port = "9164"
	}
	host, okH := os.LookupEnv("HOST")
	if !okH {
		host = "0.0.0.0" // default is for docker
	}
	fmt.Println("Listening to " + host + ":" + port + ", watching folder " + folder)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, `<!doctype html>
	<html>
		<head><title>Dir Info Exporter</title></head>
		<body>
			<h1>Dir Info Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
		</body>
	</html>`)
}

func getFolderFiles(folder string) string {
	cmd := exec.Command("sh", "-c", "find "+folder+" -type f 2>&1 | grep -v denied | wc -l")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return "0"
	}
	return strings.TrimSpace(string(stdoutStderr))
}

func getFolderSize(folder string) string {
	cmd := exec.Command("sh", "-c", "du --max-depth=0 "+folder+" 2>&1 | grep -v denied | cut -f1")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return "0"
	}
	return strings.TrimSpace(string(stdoutStderr))
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	folder, _ := os.LookupEnv("FOLDER")
	alias, okAlias := os.LookupEnv("ALIAS")
	if !okAlias {
		alias = "default"
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	res := strings.Join([]string{
		"# HELP " + alias + "_folder_info Path Information",
		"# TYPE " + alias + "_folder_info gauge",
		alias + "_folder_info{path=\"" + folder + "\"} 1",
		"# HELP " + alias + "_folder_files_count Number of files in the folder",
		"# TYPE " + alias + "_folder_files_count counter",
		alias + "_folder_files_count " + getFolderFiles(folder),
		"# HELP " + alias + "_folder_disk_size Folder Disk Size",
		"# TYPE " + alias + "_folder_disk_size counter",
		alias + "_folder_disk_size " + getFolderSize(folder),
	}, "\n")
	// fmt.Println(res)
	io.WriteString(w, res)
}
