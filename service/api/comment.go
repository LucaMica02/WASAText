package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

/* COMMENT */
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "auth token missing", http.StatusUnauthorized)
		return
	} else if auth != userIdString {
		http.Error(w, "auth token not valid", http.StatusUnauthorized)
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
	isUserConversation := isUserConversations(userId, conversationId, rt, w)
	if !isUserConversation {
		http.Error(w, "is not a user conversation", http.StatusUnauthorized)
		return
	}

	// valid request
	var comment Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// get the timestamp
	time := time.Now()
	timestamp := time.Format("2006-01-02 15:04:05")

	// create the comment
	exists, err = rt.db.CheckIfCommentExists(userId, messageId)
	if err != nil {
		http.Error(w, "Error checking if comment exists", http.StatusInternalServerError)
		return
	}
	if exists { // UPDATE
		err = rt.db.UpdateComment(userId, messageId, comment.Emoji, timestamp)
		if err != nil {
			http.Error(w, "Error updating the comment", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(comment)
		if err != nil {
			http.Error(w, "Error encoding the response", http.StatusInternalServerError)
			return
		}
	} else { // ADD
		err = rt.db.AddComment(timestamp, userId, messageId, comment.Emoji)
		if err != nil {
			http.Error(w, "Error adding the comment", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(comment)
		if err != nil {
			http.Error(w, "Error encoding the response", http.StatusInternalServerError)
			return 
		}
	}
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	// auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "auth token missing", http.StatusUnauthorized)
		return
	} else if auth != userIdString {
		http.Error(w, "auth token not valid", http.StatusUnauthorized)
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

	// valid commentId
	exists, _ = rt.db.CheckIfCommentExists(userId, messageId)
	if !exists {
		http.Error(w, "commentId not found", http.StatusNotFound)
		return
	}

	// delete the comment
	err = rt.db.DeleteComment(userId, messageId)
	if err != nil {
		http.Error(w, "Error deleting the comment", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
