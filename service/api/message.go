package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

/* MESSAGE API */
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

	// check if is a user conversation
	isUserConversation := isUserConversations(userId, conversationId, rt, w)
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

	// get the timestamp
	time := time.Now()
	timestamp := time.Format("2006-01-02 15:04:05")
	message.Timestamp = timestamp

	// create the message
	if message.RepliedTo == 0 {
		err = rt.db.CreateMessage(message.Timestamp, message.Sender, message.Conversation, message.Status, message.Type, message.Body)
		if err != nil {
			http.Error(w, "Error creating the message: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else { // is a reply
		err = rt.db.ReplyToAMessage(message.Timestamp, message.Sender, message.Conversation, message.Status, message.Type, message.Body, message.RepliedTo)
		if err != nil {
			http.Error(w, "Error creating the message: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// return the message
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
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

	// check if conversation to forward
	isUserConversation = isUserConversations(userId, resource.ResourceId, rt, w)
	if !isUserConversation {
		http.Error(w, "conversation to forward is not a user conversation", http.StatusUnauthorized)
		return
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
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}

func (rt *_router) addReceiver(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid conversationId
	conversationIdString := strings.Split(r.URL.Path, "/")[2]
	conversationId, err := strconv.Atoi(conversationIdString)
	if err != nil {
		http.Error(w, "conversationId not valid", http.StatusBadRequest)
		return
	}
	exists, _ := rt.db.CheckIfConversationExistsByConversationId(conversationId)
	if !exists {
		http.Error(w, "conversationId not found", http.StatusNotFound)
		return
	}

	// valid messageId
	messageIdString := strings.Split(r.URL.Path, "/")[4]
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

	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[6]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfUserExistsByUserId(userId)
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}

	// check if is a user conversation
	isUserConversation := isUserConversations(userId, conversationId, rt, w)
	if !isUserConversation {
		http.Error(w, "is not a user conversation", http.StatusUnauthorized)
		return
	}

	// add the receiver
	err = rt.db.AddReceiver(userId, messageId)
	if err != nil {
		http.Error(w, "Error adding the receiver", http.StatusInternalServerError)
		return
	}

	// check conversation is private
	conversation, _ := rt.db.GetConversationByConversationId(conversationId, userId)

	// count ricevitori
	receivers, _ := rt.db.GetReceivers(messageId)

	// count utenti gruppo
	if conversation.IsPrivate {
		if receivers == 2 {
			err = rt.db.UpdateMessageStatus(messageId, "received")
			if err != nil {
				http.Error(w, "Error updating the status", http.StatusInternalServerError)
				return
			}
		}
	} else {
		usersCount, _ := rt.db.GetMembersCount(conversationId)
		if receivers == usersCount {
			err = rt.db.UpdateMessageStatus(messageId, "received")
			if err != nil {
				http.Error(w, "Error updating the status", http.StatusInternalServerError)
				return
			}
		}
	}

	// return the userId
	var resourceId ResourceId
	resourceId.ResourceId = userId
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resourceId)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}

func (rt *_router) addReader(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid conversationId
	conversationIdString := strings.Split(r.URL.Path, "/")[2]
	conversationId, err := strconv.Atoi(conversationIdString)
	if err != nil {
		http.Error(w, "conversationId not valid", http.StatusBadRequest)
		return
	}
	exists, _ := rt.db.CheckIfConversationExistsByConversationId(conversationId)
	if !exists {
		http.Error(w, "conversationId not found", http.StatusNotFound)
		return
	}

	// valid messageId
	messageIdString := strings.Split(r.URL.Path, "/")[4]
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

	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[6]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfUserExistsByUserId(userId)
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}

	// check if is a user conversation
	isUserConversation := isUserConversations(userId, conversationId, rt, w)
	if !isUserConversation {
		http.Error(w, "is not a user conversation", http.StatusUnauthorized)
		return
	}

	// add the reader
	err = rt.db.AddReader(userId, messageId)
	if err != nil {
		http.Error(w, "Error adding the receiver", http.StatusInternalServerError)
		return
	}

	// check conversation is private
	conversation, _ := rt.db.GetConversationByConversationId(conversationId, userId)

	// count ricevitori
	readers, _ := rt.db.GetReaders(messageId)

	// count utenti gruppo
	if conversation.IsPrivate {
		if readers == 2 {
			err = rt.db.UpdateMessageStatus(messageId, "read")
			if err != nil {
				http.Error(w, "Error updating the status", http.StatusInternalServerError)
				return
			}
		}
	} else {
		usersCount, _ := rt.db.GetMembersCount(conversationId)
		if readers == usersCount {
			err = rt.db.UpdateMessageStatus(messageId, "read")
			if err != nil {
				http.Error(w, "Error updating the status", http.StatusInternalServerError)
				return
			}
		}
	}

	// return the userId
	var resourceId ResourceId
	resourceId.ResourceId = userId
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resourceId)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}
