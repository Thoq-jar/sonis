package main

import (
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sonis/net"
	"sonis/utils"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_ = mime.AddExtensionType(".flac", "audio/flac")

	if err := utils.BuildLibrary("music"); err != nil {
		log.Fatalf("failed to build library: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/library", net.HandleLibrary)
	mux.HandleFunc("/api/rescan", net.HandleRescan)
	mux.HandleFunc("/stream/", net.HandleStream)

	uiDir := "web"
	buildIndex := filepath.Join("web", "build", "index.html")

	if st, err := os.Stat(buildIndex); err == nil && !st.IsDir() {
		uiDir = filepath.Join("web", "build")
	}

	fs := http.FileServer(http.Dir(uiDir))
	mux.Handle("/", fs)

	addr := ":8080"

	log.Printf("sonis server listening on %s", addr)
	log.Printf("open http://localhost%s/ in your browser", addr)

	if err := http.ListenAndServe(addr, net.WithCommonHeaders(net.LoggingMiddleware(mux))); err != nil {
		log.Fatal(err)
	}
}
