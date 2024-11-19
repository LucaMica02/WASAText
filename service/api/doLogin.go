package api

import (
	"encoding/json"
	"net/http"
)

func doLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only Post Allowed", http.StatusMethodNotAllowed)
		return
	}

	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Bad Request, Invalide json", http.StatusBadRequest)
		return
	}

	if request.username == "" {
		http.Error(w, "Bad Request, Username required", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	// Log in or Create new user
}
