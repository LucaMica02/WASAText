package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Login route
	rt.router.POST("/session", rt.doLogin)

	// User API routes
	rt.router.GET("/users", rt.getUsers)
	rt.router.GET("/users/:userId", rt.getUser)
	rt.router.PUT("/users/:userId/username", rt.setMyUserName)
	rt.router.PUT("/users/:userId/photo", rt.setMyPhoto)

	// Conversation API routes
	rt.router.GET("/users/:userId/conversations", rt.getMyConversations)
	rt.router.GET("/users/:userId/conversations/:conversationsId", rt.getConversation)
	rt.router.POST("/users/:userId/conversations", rt.createConversation)

	// Message API routes
	rt.router.POST("/users/:userId/conversations/:conversationsId/messages", rt.sendMessage)
	rt.router.DELETE("/users/:userId/conversations/:conversationsId/messages/:messageId", rt.deleteMessage)
	rt.router.POST("/users/:userId/conversations/:conversationsId/messages/:messageId/forward", rt.forwardMessage)
	rt.router.PUT("/conversations/:conversationsId/messages/:messageId/status", rt.updateMessageStatus)
	rt.router.PUT("/conversations/:conversationsId/messages/:messageId/receivers/:userIdDest", rt.addReceiver)
	rt.router.PUT("/conversations/:conversationsId/messages/:messageId/readers/:userIdDest", rt.addReader)

	// Comment API routes
	rt.router.PUT("/users/:userId/conversations/:conversationsId/messages/:messageId/comment", rt.commentMessage)
	rt.router.DELETE("/users/:userId/conversations/:conversationsId/messages/:messageId/comment/:commentId", rt.uncommentMessage)

	// Group API routes
	rt.router.POST("/users/:userId/groups", rt.createGroup)
	rt.router.PUT("/users/:userId/groups/:groupId/members", rt.addToGroup)
	rt.router.DELETE("/users/:userId/groups/:groupId/members", rt.leaveGroup)
	rt.router.PUT("/users/:userId/groups/:groupId/name", rt.setGroupName)
	rt.router.PUT("/users/:userId/groups/:groupId/description", rt.setGroupDescription)
	rt.router.PUT("/users/:userId/groups/:groupId/photo", rt.setGroupPhoto)

	// Get Image route
	rt.router.GET("/images", rt.getPhotoHandler)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
