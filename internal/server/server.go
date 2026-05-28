package server

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/why19970628/schemaforge/internal/convert"
)

//go:embed web/dist/*
var webFS embed.FS

func ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	Register(mux)
	return http.ListenAndServe(addr, mux)
}

func Register(mux *http.ServeMux) {
	staticFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	mux.HandleFunc("/api/convert", handleConvert)
	mux.HandleFunc("/api/modes", handleModes)
}

func handleModes(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"modes": convert.Modes()})
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}
	defer r.Body.Close()
	var req convert.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	resp, err := convert.Convert(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
