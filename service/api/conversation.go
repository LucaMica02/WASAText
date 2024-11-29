package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

/* CONVERSATIONS API */
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
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

	// return the conversations
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(getUserConversations(userId, rt, w))
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
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
	isUserConversation := isUserConversations(userId, conversationId, rt, w)
	if !isUserConversation {
		http.Error(w, "is not a user conversation", http.StatusUnauthorized)
		return
	}

	// return the conversation
	conversations, _ := rt.db.GetConversationByConversationId(conversationId, userId)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
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
		conversationId, err = rt.db.CreatePrivateConversation(users[0].ResourceId, users[1].ResourceId)
		if err != nil {
			http.Error(w, "Error creating the conversation", http.StatusBadRequest)
		}
	} else {
		conversationId, err = rt.db.CreatePrivateConversation(users[1].ResourceId, users[0].ResourceId)
		if err != nil {
			http.Error(w, "Error creating the conversation", http.StatusBadRequest)
		}
	}

	// return the conversation
	conversations, _ := rt.db.GetConversationByConversationId(conversationId, userId)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}
