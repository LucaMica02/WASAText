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
	"time"

	"github.com/julienschmidt/httprouter"
)

// FUNCTION FOR API ROUTES

/* LOGIN */
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

/* USERS */
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

/* IMAGE HANDLER */
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

/* CONVERSATIONS */
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
	// valid userId
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

	// valid conversationId
	conversationIdString := strings.Split(r.URL.Path, "/")[4]
	conversationId, err := strconv.Atoi(conversationIdString)
	if err != nil {
		http.Error(w, "conversationId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfConversationExistsByConversationId(conversationId)
	if !exists {
		http.Error(w, "conversationId not found", http.StatusNotFound)
		return
	}

	// check if is a user conversation
	isUserConversation := false
	groupConversations, _ := rt.db.GetGroupConversationsByUserId(userId)
	privateConversations, _ := rt.db.GetPrivateConversationsByUserId(userId)
	totalConversations := append(groupConversations, privateConversations...)
	for _, conversation := range totalConversations {
		if conversation.ResourceId == conversationId {
			isUserConversation = true
		}
	}
	if !isUserConversation {
		http.Error(w, "is not a user conversation", http.StatusUnauthorized)
		return
	}

	// return the conversation
	conversations, _ := rt.db.GetConversationByConversationId(conversationId, userId)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conversations)
}

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
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
	// valid request
	var users []ResourceId
	err = json.NewDecoder(r.Body).Decode(&users)
	if err != nil || len(users) != 2 {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	if users[0].ResourceId != userId && users[1].ResourceId != userId {
		http.Error(w, "Invalid Request", http.StatusUnauthorized)
		return
	}
	// constraint userId_1 < userId_2
	var conversationId int
	if users[0].ResourceId < users[1].ResourceId {
		conversationId, _ = rt.db.CreatePrivateConversation(users[0].ResourceId, users[1].ResourceId)
	} else {
		conversationId, _ = rt.db.CreatePrivateConversation(users[1].ResourceId, users[0].ResourceId)
	}
	conversations, _ := rt.db.GetConversationByConversationId(conversationId, userId)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(conversations)
	// API & DEBUG
}

/* MESSAGE */
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
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

	// valid conversationId
	conversationIdString := strings.Split(r.URL.Path, "/")[4]
	conversationId, err := strconv.Atoi(conversationIdString)
	if err != nil {
		http.Error(w, "conversationId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfConversationExistsByConversationId(conversationId)
	if !exists {
		http.Error(w, "conversationId not found", http.StatusNotFound)
		return
	}

	// check if is a user conversation
	isUserConversation := false
	groupConversations, _ := rt.db.GetGroupConversationsByUserId(userId)
	privateConversations, _ := rt.db.GetPrivateConversationsByUserId(userId)
	totalConversations := append(groupConversations, privateConversations...)
	for _, conversation := range totalConversations {
		if conversation.ResourceId == conversationId {
			isUserConversation = true
		}
	}
	if !isUserConversation {
		http.Error(w, "is not a user conversation", http.StatusUnauthorized)
		return
	}

	// valid request
	var message Message
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	message.Sender = userId
	message.Conversation = conversationId
	message.Status = "delivered"

	// create the message
	err = rt.db.CreateMessage(message.Timestamp, message.Sender, message.Conversation, message.Status, message.Type, message.Body, message.RepliedTo, message.ForwardedFrom)
	if err != nil {
		http.Error(w, "Error creating the message", http.StatusInternalServerError)
		return
	}

	// return the message
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
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

	// valid conversationId
	conversationIdString := strings.Split(r.URL.Path, "/")[4]
	conversationId, err := strconv.Atoi(conversationIdString)
	if err != nil {
		http.Error(w, "conversationId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfConversationExistsByConversationId(conversationId)
	if !exists {
		http.Error(w, "conversationId not found", http.StatusNotFound)
		return
	}

	// valid messageId
	messageIdString := strings.Split(r.URL.Path, "/")[6]
	messageId, err := strconv.Atoi(messageIdString)
	if err != nil {
		http.Error(w, "messageId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfMessageExistsByMessageId(messageId)
	if !exists {
		http.Error(w, "messageId not found", http.StatusNotFound)
		return
	}

	// check if is a user conversation
	isUserConversation := false
	groupConversations, _ := rt.db.GetGroupConversationsByUserId(userId)
	privateConversations, _ := rt.db.GetPrivateConversationsByUserId(userId)
	totalConversations := append(groupConversations, privateConversations...)
	for _, conversation := range totalConversations {
		if conversation.ResourceId == conversationId {
			isUserConversation = true
		}
	}
	if !isUserConversation {
		http.Error(w, "is not a user conversation", http.StatusUnauthorized)
		return
	}

	// delete the message
	err = rt.db.DeleteMessage(messageId)
	if err != nil {
		http.Error(w, "Error deleting the message", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
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

	// valid conversationId
	conversationIdString := strings.Split(r.URL.Path, "/")[4]
	conversationId, err := strconv.Atoi(conversationIdString)
	if err != nil {
		http.Error(w, "conversationId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfConversationExistsByConversationId(conversationId)
	if !exists {
		http.Error(w, "conversationId not found", http.StatusNotFound)
		return
	}

	// valid messageId
	messageIdString := strings.Split(r.URL.Path, "/")[6]
	messageId, err := strconv.Atoi(messageIdString)
	if err != nil {
		http.Error(w, "messageId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfMessageExistsByMessageId(messageId)
	if !exists {
		http.Error(w, "messageId not found", http.StatusNotFound)
		return
	}

	// check if is a user conversation
	isUserConversation := false
	groupConversations, _ := rt.db.GetGroupConversationsByUserId(userId)
	privateConversations, _ := rt.db.GetPrivateConversationsByUserId(userId)
	totalConversations := append(groupConversations, privateConversations...)
	for _, conversation := range totalConversations {
		if conversation.ResourceId == conversationId {
			isUserConversation = true
		}
	}
	if !isUserConversation {
		http.Error(w, "is not a user conversation", http.StatusUnauthorized)
		return
	}

	// valid request
	var resource ResourceId
	err = json.NewDecoder(r.Body).Decode(&resource)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfConversationExistsByConversationId(resource.ResourceId)
	if !exists {
		http.Error(w, "conversationId to forward not found", http.StatusNotFound)
		return
	}
	for _, conversation := range totalConversations {
		if conversation.ResourceId == resource.ResourceId {
			isUserConversation = true
		}
	}
	if !isUserConversation {
		http.Error(w, "conversation to forward is not a user conversation", http.StatusUnauthorized)
		return
	}

	// get the timestamp
	time := time.Now()
	timestamp := time.Format("2006-01-02 15:04:05")

	// create the message
	message, err := rt.db.ForwardMessage(messageId, userId, resource.ResourceId, timestamp)
	if err != nil {
		http.Error(w, "Error creating the message", http.StatusInternalServerError)
		return
	}

	// return the message
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}
