package net

import (
	"net/http"
	"sonis/utils"
	"strings"
)

func HandleStream(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/stream/")

	if id == "" || id == r.URL.Path {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	t, ok := utils.TrackByID[id]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "audio/flac")
	http.ServeFile(w, r, t.Path)
}
