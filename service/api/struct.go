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
