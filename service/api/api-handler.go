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
	rt.router.POST("/conversations", rt.createConversation)

	// Get Image route
	rt.router.GET("/images", rt.getPhotoHandler)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
