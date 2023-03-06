package handlers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"product-images/pkg/files"
	"strconv"

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

func (f *Files) MultiPartUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024)
	if err != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(w, "Expected multipart form data", http.StatusBadRequest)
		return
	}
	// getting values from the UI form
	// need to validate ID is an int
	id, idErr := strconv.Atoi(r.FormValue("id"))
	f.log.Info("form id", "id", id)
	if idErr != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(w, "Expected integer id", http.StatusBadRequest)
		return
	}

	// getting file from form request
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(w, "Expected file", http.StatusBadRequest)
		return
	}
	f.saveFile(fmt.Sprint(id), fileHeader.Filename, w, file)
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
