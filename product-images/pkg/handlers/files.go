package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"product-images/pkg/files"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

type Files struct {
	log   hclog.Logger
	store files.Storage
}

func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{
		log:   l,
		store: s,
	}
}

func (f *Files) Upload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", filename)
	f.saveFile(id, filename, w, r.Body)
}

func (f *Files) saveFile(id string, path string, w http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file for image", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
	}
}
