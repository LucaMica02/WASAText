package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
)

func getUserConversations(userId int, rt *_router, w http.ResponseWriter) []database.ConversationId {
	// get the conversations
	groupConversations, err := rt.db.GetGroupConversationsByUserId(userId)
	if err != nil {
		http.Error(w, "Error getting the group conversations", http.StatusInternalServerError)
		return nil
	}
	privateConversations, err := rt.db.GetPrivateConversationsByUserId(userId)
	if err != nil {
		http.Error(w, "Error getting the private conversations", http.StatusInternalServerError)
		return nil
	}
	totalConversations := append(groupConversations, privateConversations...)
	return totalConversations
}

func isUserConversations(userId int, conversationId int, rt *_router, w http.ResponseWriter) bool {
	// return true if the conversationId is a user conversation false otherwise
	isUserConversation := false
	totalConversations := getUserConversations(userId, rt, w)
	for _, conversation := range totalConversations {
		if conversation.ResourceId == conversationId {
			isUserConversation = true
		}
	}
	return isUserConversation
}
