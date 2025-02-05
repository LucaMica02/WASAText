package database

type ConversationId struct {
	ResourceId int `json:"resourceId"`
}

type Conversation struct {
	Name        string    `json:"conversationName"`
	Description string    `json:"description"`
	Messages    []Message `json:"messages"`
	IsPrivate   bool      `json:"isPrivate"`
	PhotoUrl    string    `json:"photoUrl"`
}

type Message struct {
	ResourceId    int    `json:"resourceId"`
	Timestamp     string `json:"timestamp"`
	Sender        int    `json:"sender"`
	Conversation  int    `json:"conversation"`
	Status        string `json:"status"`
	Type          string `json:"type"`
	Body          string `json:"body"`
	RepliedTo     int    `json:"repliedTo"`
	ForwardedFrom int    `json:"forwardedFrom"`
	Comments      int    `json:"comments"`
}

type User struct {
	ResourceId int    `json:"resourceId"`
	Username   string `json:"username"`
	PhotoUrl   string `json:"PhotoUrl"`
}
