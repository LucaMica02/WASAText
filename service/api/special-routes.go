package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

/* LOGIN */
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
	}

	w.Header().Set("content-type", "application/json")
	exists, err := rt.db.CheckIfUserExistsByUsername(loginRequest.Username)
	if err != nil {
		http.Error(w, "Error checking if the user exists", http.StatusInternalServerError)
	}
	if exists {
		// return userId
		id, _ := rt.db.GetUserId(loginRequest.Username)
		response := LoginResponse{ResourceId: id}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error encoding the response", http.StatusInternalServerError)
		}
	} else {
		// create new user and return userId
		_ = rt.db.CreateUser(loginRequest.Username)
		id, _ := rt.db.GetUserId(loginRequest.Username)
		response := LoginResponse{ResourceId: id}
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error encoding the response", http.StatusInternalServerError)
		}
	}
}

/* IMAGE HANDLER */
func (rt *_router) getPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid the path
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "No path provided", http.StatusBadRequest)
		return
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// serve the file
	ext := filepath.Ext(path)
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("content-type", "image/jpeg")
	case ".png":
		w.Header().Set("content-type", "image/png")
	}
	http.ServeFile(w, r, path)
}

/* SAVE PHOTO */
func (rt *_router) savePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid the path
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "No path provided", http.StatusBadRequest)
		return
	}

	// parse the file
	err := r.ParseMultipartForm(10 << 25) // max 25 MB
	if err != nil {
		http.Error(w, "Error parsing form-data", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Error loading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// create directory if not exists
	imagesDir := "images"
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		err = os.Mkdir(imagesDir, os.ModePerm)
		if err != nil {
			http.Error(w, "Error creating images directory", http.StatusInternalServerError)
			return
		}
	}

	// create the path
	ext := filepath.Ext(header.Filename)
	if ext != ".jpeg" && ext != ".png" {
		http.Error(w, "Only .jpeg and .png image allowed", http.StatusBadRequest)
		return
	}
	filePath := filepath.Join(imagesDir, path+ext)

	// write the file
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
	}
	defer out.Close()

	// copy the content into the destination file
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
	}

	// return the User
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"path": filePath}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}
