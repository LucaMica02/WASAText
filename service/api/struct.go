package api

import "time"

type MessageStatus string

const (
	Delivered MessageStatus = "delivered"
	Received  MessageStatus = "received"
	Read      MessageStatus = "read"
)

type MessageType string

const (
	Image MessageType = "image"
	Text  MessageType = "text"
)

type LoginRequest struct {
	username string
}

type LoginResponse struct {
	userId int
}

type User struct {
	userId   int
	username string
	photoUrl string
}

type PrivateConversation struct {
	ConversationId int
	UserId1        int
	UserId2        int
}

type Group struct {
	groupId        int
	name           string
	description    string
	photoUrl       string
	users          []int
	ConversationId int
}

type Message struct {
	messageId     int
	time          time.Time
	sender        int
	conversation  int
	receivers     []int
	readers       []int
	messageStatus MessageStatus
	messageType   MessageType
	content       string
}

type Comment struct {
	commentId int
	time      time.Time
	sender    int
	message   int
	reaction  string
}
