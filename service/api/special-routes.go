package api

import (
	"encoding/json"
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
		return 
	}

	w.Header().Set("content-type", "application/json")
	exists, err := rt.db.CheckIfUserExistsByUsername(loginRequest.Username)
	if err != nil {
		http.Error(w, "Error checking if the user exists", http.StatusInternalServerError)
		return 
	}
	if exists {
		// return userId
		id, _ := rt.db.GetUserId(loginRequest.Username)
		response := LoginResponse{ResourceId: id}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error encoding the response", http.StatusInternalServerError)
			return 
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
			return 
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
