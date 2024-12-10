package database

type ConversationId struct {
	ResourceId int `json:"resourceId"`
}

type Conversation struct {
	Name     string    `json:"conversationName"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Timestamp     string `json:"timestamp"`
	Sender        int    `json:"sender"`
	Conversation  int    `json:"conversation"`
	Status        string `json:"status"`
	Type          string `json:"type"`
	Body          string `json:"body"`
	RepliedTo     int    `json:"repliedTo"`
	ForwardedFrom int    `json:"forwardedFrom"`
}

type User struct {
	ResourceId int  `json:"resourceId"`
	Username string `json:"username"`
	PhotoUrl string `json:"PhotoUrl"`
}
