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

// FUNCTION FOR API ROUTES
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
	}

	w.Header().Set("content-type", "application/json")
	exists, _ := rt.db.CheckIfUserExistsByUsername(loginRequest.Username)
	if exists {
		// return userId
		id, _ := rt.db.GetUserId(loginRequest.Username)
		response := LoginResponse{ResourceId: id}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		// create new user and return userId
		_ = rt.db.CreateUser(loginRequest.Username)
		id, _ := rt.db.GetUserId(loginRequest.Username)
		response := LoginResponse{ResourceId: id}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	username := r.URL.Query().Get("username")
	if username == "" {
		// search for all the users
		users, _ := rt.db.GetAllUsers()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(users)
	} else {
		// search for the specified user
		exists, _ := rt.db.CheckIfUserExistsByUsername(username)
		if exists { // return the User
			user, _ := rt.db.GetUserByUsername(username)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
		} else {
			http.Error(w, "User not found", http.StatusNotFound)
		}
	}
}

func (rt *_router) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	user, err := rt.db.GetUserById(userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}
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
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIdString := strings.Split(r.URL.Path, "/")[2]
	err := r.ParseMultipartForm(10 << 25) // max 25 MB
	if err != nil {
		http.Error(w, "Error parsing form-data", http.StatusBadRequest)
		return
	}
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
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

	err = rt.db.UpdatePhotoUrl(filePath, userId)
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
	json.NewEncoder(w).Encode(user)
}

func (rt *_router) getPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "No path provided", http.StatusBadRequest)
		return
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("content-type", "image/jpeg")
	case ".png":
		w.Header().Set("content-type", "image/png")
	}
	http.ServeFile(w, r, path)
}

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}
	exists, _ := rt.db.CheckIfUserExistsByUserId(userId)
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}
	groupConversations, _ := rt.db.GetGroupConversationsByUserId(userId)
	privateConversations, _ := rt.db.GetPrivateConversationsByUserId(userId)
	totalConversations := append(groupConversations, privateConversations...)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(totalConversations)
}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}
	exists, _ := rt.db.CheckIfUserExistsByUserId(userId)
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}
	conversationIdString := strings.Split(r.URL.Path, "/")[2]
	conversationId, err := strconv.Atoi(conversationIdString)
	if err != nil {
		http.Error(w, "conversationId not valid", http.StatusBadRequest)
		return
	}
	conversation, _ := rt.db.GetConversationByConversationId(conversationId)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conversation)
	// security : exists and is user's chat
	// debug Conversation
}

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
