package api

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	ResourceId int `json:"resourceId"`
}

type Username struct {
	Username string `json:"username"`
}

type ResourceId struct {
	ResourceId int `json:"resourceId"`
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
