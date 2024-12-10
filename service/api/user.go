package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

/* USERS API */
func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	username := r.URL.Query().Get("username")
	if username == "" {
		// search for all the users
		users, err := rt.db.GetAllUsers()
		if err != nil {
			http.Error(w, "Error getting the users", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			http.Error(w, "Error encoding the response", http.StatusInternalServerError)
		}
	} else {
		// search for the specified user
		exists, err := rt.db.CheckIfUserExistsByUsername(username)
		if err != nil {
			http.Error(w, "Error checking the user", http.StatusInternalServerError)
		}
		if exists { // return the User
			user, err := rt.db.GetUserByUsername(username)
			if err != nil {
				http.Error(w, "Error getting the user", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(user)
			if err != nil {
				http.Error(w, "Error encoding the response", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "User not found", http.StatusNotFound)
		}
	}
}

func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}

	// auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "auth token missing", http.StatusUnauthorized)
		return
	}

	// check if userId exists
	exists, err := rt.db.CheckIfUserExistsByUserId(userId)
	if err != nil {
		http.Error(w, "Error checking the user", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}

	// get the user
	user, err := rt.db.GetUserById(userId)
	if err != nil {
		http.Error(w, "Error getting the user", http.StatusInternalServerError)
		return
	}

	// return the user
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}

	// auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "auth token missing", http.StatusUnauthorized)
		return
	} else if auth != userIdString {
		http.Error(w, "auth token not valid", http.StatusUnauthorized)
		return
	}

	// check if userId exists
	exists, err := rt.db.CheckIfUserExistsByUserId(userId)
	if err != nil {
		http.Error(w, "Error checking the user", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}

	// update the username
	var username Username
	err = json.NewDecoder(r.Body).Decode(&username)
	if err != nil {
		http.Error(w, "username not valid", http.StatusBadRequest)
		return
	}
	err = rt.db.UpdateUsername(username.Username, userId)
	if err != nil {
		http.Error(w, "username not unique", http.StatusBadRequest)
		return
	}
	user, err := rt.db.GetUserById(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// return the user
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}

	// auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "auth token missing", http.StatusUnauthorized)
		return
	} else if auth != userIdString {
		http.Error(w, "auth token not valid", http.StatusUnauthorized)
		return
	}

	// check if userId exists
	exists, err := rt.db.CheckIfUserExistsByUserId(userId)
	if err != nil {
		http.Error(w, "Error checking the user", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}

	// parse the file
	err = r.ParseMultipartForm(10 << 25) // max 25 MB
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
	fileName := fmt.Sprintf("user_%s_photo%s", userIdString, ext)
	filePath := filepath.Join(imagesDir, fileName)

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

	// update the user's photoUrl
	err = rt.db.UpdateUserPhotoUrl(filePath, userId)
	if err != nil {
		http.Error(w, "Error saving the photoUrl in the DB", http.StatusInternalServerError)
		return
	}

	// return the User
	user, err := rt.db.GetUserById(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}
