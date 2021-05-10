package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

const (
	DEFAULT_PORT = "8888"
)

func main() {
	var (
		port string
	)
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}
	log.SetOutput(os.Stdout)
	log.Printf("Starting on port %s", port)
	http.HandleFunc("/", handler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/download/1G", oneGbDownloadHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "text/html")
	hostname, _ := os.Hostname()
	responseHTML := `<!DOCTYPE html>
	<html>
	<head>
	  <meta charset='UTF-8'>
	  <title>%s</title>
	  <link rel="icon" href="data:,">
	</head>
	<body>
	  <h2>Server Info</h2>
	  <pre>%s</pre>
	  <h2>Request Info</h2>
	  <pre>%s</pre>
	  <h3>Header</h3>
	  <pre>%s</pre>
	  <h2>other</h2>
	  <a href='/download'>file download</a>
	</body>
	</html>
	`
	headers := make([]string, len(r.Header))
	index := 0
	for k, _ := range r.Header {
		headers[index] = k
		index++
	}
	sort.Strings(headers)
	servInfo := fmt.Sprintf("  Hostname: %s", hostname)
	reqInfo := fmt.Sprintf("  Client:   %s\n  URL:      %s", r.RemoteAddr, r.URL)
	headerInfo := fmt.Sprintf("  Host: %s\n", r.Host)
	for _, key := range headers {
		headerInfo += fmt.Sprintf("  %s: %s\n", key, r.Header[key][0])
	}
	log.Printf("Client: %s, Method: %s, URL: %s, Duration: %s", r.RemoteAddr, r.Method, r.URL.EscapedPath(), time.Since(start))
	log.Printf("Request Headers:\n%s", headerInfo)
	w.Write([]byte(fmt.Sprintf(responseHTML, hostname, servInfo, reqInfo, headerInfo)))
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Client: %s, Method: %s, URL: %s", r.RemoteAddr, r.Method, r.URL.EscapedPath())
	w.Header().Set("Content-Type", "text/html")
	hostname, _ := os.Hostname()
	responseHTML := `<!DOCTYPE html>
	<html>
	<head>
	  <meta charset='UTF-8'>
	  <title>%s</title>
	  <link rel="icon" href="data:,">
	</head>
	<body>
	  <ul>
	  <li><a href='/download/1G'>1G</a></li>
	  </ul>
	</body>
	</html>
	`
	w.Write([]byte(fmt.Sprintf(responseHTML, hostname)))
}

func oneGbDownloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, URL: %s", r.Method, r.URL.EscapedPath())
	//w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\"dummy\"")
	dummyFile := "/tmp/1g"
	_, err := os.Stat(dummyFile)
	if os.IsNotExist(err) {
		f, _ := os.Create(dummyFile)
		f.Truncate(1e9)
	}
	http.ServeFile(w, r, dummyFile)
}
