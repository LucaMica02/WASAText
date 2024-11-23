package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getHelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Hello World!"))
}

func (rt *_router) getText(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name, err := rt.db.GetName()
	w.Header().Set("content-type", "text/plain")
	if err != nil {
		_, _ = w.Write([]byte("No Text!"))
	} else {
		_, _ = w.Write([]byte(name))
	}
}

func (rt *_router) addText(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := r.URL.Query().Get("text")
	rt.db.SetName(name)
}

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	ResourceId int `json:"resourceId"`
}

// FUNCTION FOR API ROUTES
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
	}

	exists, _ := rt.db.CheckIfUserExists(loginRequest.Username)
	if exists {
		// return userId
		id, _ := rt.db.GetUserId(loginRequest.Username)
		response := LoginResponse{ResourceId: id}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	} else {
		// create new user and return userId
		_ = rt.db.CreateUser(loginRequest.Username)
		id, _ := rt.db.GetUserId(loginRequest.Username)
		response := LoginResponse{ResourceId: id}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
