package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

/* GROUP */
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}
	exists, _ := rt.db.CheckIfUserExistsByUserId(userId)
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}

	// auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "auth token missing", http.StatusUnauthorized)
		return
	} else if auth != userIdString {
		http.Error(w, "auth token not valid", http.StatusUnauthorized)
		return
	}

	// valid request
	var group Group
	err = json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// check if user is a partecipant
	isUserInPartecipants := false
	for _, user := range group.Partecipants {
		if user.ResourceId == userId {
			isUserInPartecipants = true
		}
	}
	if !isUserInPartecipants {
		http.Error(w, "The user have to be a group partecipant", http.StatusBadRequest)
		return
	}

	// create the group
	groupId, err := rt.db.CreateGroupConversation(group.Name, group.Description, group.Photo)
	if err != nil {
		http.Error(w, "Error creating the group", http.StatusInternalServerError)
		return
	}

	// add partecipants in UserGroup
	for _, user := range group.Partecipants {
		err = rt.db.AddUserToGroup(user.ResourceId, groupId)
		if err != nil {
			http.Error(w, "Error creating UserGroup entry", http.StatusInternalServerError)
			return
		}
	}

	// return the group
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(group)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}
	exists, _ := rt.db.CheckIfUserExistsByUserId(userId)
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}

	// auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "auth token missing", http.StatusUnauthorized)
		return
	} else if auth != userIdString {
		http.Error(w, "auth token not valid", http.StatusUnauthorized)
		return
	}

	// valid groupId
	groupIdString := strings.Split(r.URL.Path, "/")[4]
	groupId, err := strconv.Atoi(groupIdString)
	if err != nil {
		http.Error(w, "groupId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfGroupExistsByGroupId(groupId)
	if !exists {
		http.Error(w, "groupId not found", http.StatusNotFound)
		return
	}

	// valid request
	userIdToAddString := r.URL.Query().Get("userId")
	userIdToAdd, err := strconv.Atoi(userIdToAddString)
	if err != nil {
		http.Error(w, "userId to add not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfUserExistsByUserId(userIdToAdd)
	if !exists {
		http.Error(w, "userId to add not found", http.StatusNotFound)
		return
	}

	// userId is group partecipants and userId to add is not
	userIdIsPartecipant, err := rt.db.CheckIfUserIsPartecipant(userId, groupId)
	if err != nil {
		http.Error(w, "Error checking userId is partecipant", http.StatusInternalServerError)
		return
	}
	userIdToAddIsPartecipant, err := rt.db.CheckIfUserIsPartecipant(userIdToAdd, groupId)
	if err != nil {
		http.Error(w, "Error checking userId to add is partecipant", http.StatusInternalServerError)
		return
	}
	if !userIdIsPartecipant {
		http.Error(w, "User Id is not a group partecipant", http.StatusUnauthorized)
		return
	}
	if userIdToAddIsPartecipant {
		http.Error(w, "User id to add already a group partecipant", http.StatusBadRequest)
		return
	}

	// add partecipant to UserGroup
	err = rt.db.AddUserToGroup(userIdToAdd, groupId)
	if err != nil {
		http.Error(w, "Error creating UserGroup entry", http.StatusInternalServerError)
		return
	}

	// return the resourceId
	var resourceId ResourceId
	resourceId.ResourceId = userIdToAdd
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(resourceId)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}
	exists, _ := rt.db.CheckIfUserExistsByUserId(userId)
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}

	// auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "auth token missing", http.StatusUnauthorized)
		return
	} else if auth != userIdString {
		http.Error(w, "auth token not valid", http.StatusUnauthorized)
		return
	}

	// valid groupId
	groupIdString := strings.Split(r.URL.Path, "/")[4]
	groupId, err := strconv.Atoi(groupIdString)
	if err != nil {
		http.Error(w, "groupId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfGroupExistsByGroupId(groupId)
	if !exists {
		http.Error(w, "groupId not found", http.StatusNotFound)
		return
	}

	// userId is group partecipants
	userIdIsPartecipant, err := rt.db.CheckIfUserIsPartecipant(userId, groupId)
	if err != nil {
		http.Error(w, "Error checking userId is partecipant", http.StatusInternalServerError)
		return
	}
	if !userIdIsPartecipant {
		http.Error(w, "User Id is not a group partecipant", http.StatusUnauthorized)
		return
	}

	// remove the user from UserGroup
	err = rt.db.LeaveGroup(userId, groupId)
	if err != nil {
		http.Error(w, "Error deleting UserGroup entry", http.StatusInternalServerError)
		return
	}

	// return the resourceId
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}
	exists, _ := rt.db.CheckIfUserExistsByUserId(userId)
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}

	// auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "auth token missing", http.StatusUnauthorized)
		return
	} else if auth != userIdString {
		http.Error(w, "auth token not valid", http.StatusUnauthorized)
		return
	}

	// valid groupId
	groupIdString := strings.Split(r.URL.Path, "/")[4]
	groupId, err := strconv.Atoi(groupIdString)
	if err != nil {
		http.Error(w, "groupId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfGroupExistsByGroupId(groupId)
	if !exists {
		http.Error(w, "groupId not found", http.StatusNotFound)
		return
	}

	// userId is group partecipants
	userIdIsPartecipant, err := rt.db.CheckIfUserIsPartecipant(userId, groupId)
	if err != nil {
		http.Error(w, "Error checking userId is partecipant", http.StatusInternalServerError)
		return
	}
	if !userIdIsPartecipant {
		http.Error(w, "User Id is not a group partecipant", http.StatusUnauthorized)
		return
	}

	// valid request
	var groupName GroupName
	err = json.NewDecoder(r.Body).Decode(&groupName)
	length := len(groupName.GroupName)
	if err != nil || length < 2 || length > 20 {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// update the group name
	err = rt.db.UpdateGroupName(groupId, groupName.GroupName)
	if err != nil {
		http.Error(w, "Error updating group name", http.StatusInternalServerError)
		return
	}

	// return the resourceId
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(groupName)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}

func (rt *_router) setGroupDescription(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}
	exists, _ := rt.db.CheckIfUserExistsByUserId(userId)
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}

	// auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "auth token missing", http.StatusUnauthorized)
		return
	} else if auth != userIdString {
		http.Error(w, "auth token not valid", http.StatusUnauthorized)
		return
	}

	// valid groupId
	groupIdString := strings.Split(r.URL.Path, "/")[4]
	groupId, err := strconv.Atoi(groupIdString)
	if err != nil {
		http.Error(w, "groupId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfGroupExistsByGroupId(groupId)
	if !exists {
		http.Error(w, "groupId not found", http.StatusNotFound)
		return
	}

	// userId is group partecipants
	userIdIsPartecipant, err := rt.db.CheckIfUserIsPartecipant(userId, groupId)
	if err != nil {
		http.Error(w, "Error checking userId is partecipant", http.StatusInternalServerError)
		return
	}
	if !userIdIsPartecipant {
		http.Error(w, "User Id is not a group partecipant", http.StatusUnauthorized)
		return
	}

	// valid request
	var groupDescription GroupDescription
	err = json.NewDecoder(r.Body).Decode(&groupDescription)
	length := len(groupDescription.GroupDescription)
	if err != nil || length < 1 || length > 150 {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// update the group description
	err = rt.db.UpdateGroupDescription(groupId, groupDescription.GroupDescription)
	if err != nil {
		http.Error(w, "Error updating group description", http.StatusInternalServerError)
		return
	}

	// return the resourceId
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(groupDescription)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// valid userId
	userIdString := strings.Split(r.URL.Path, "/")[2]
	userId, err := strconv.Atoi(userIdString)
	if err != nil {
		http.Error(w, "userId not valid", http.StatusBadRequest)
		return
	}
	exists, _ := rt.db.CheckIfUserExistsByUserId(userId)
	if !exists {
		http.Error(w, "userId not found", http.StatusNotFound)
		return
	}

	// auth
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "auth token missing", http.StatusUnauthorized)
		return
	} else if auth != userIdString {
		http.Error(w, "auth token not valid", http.StatusUnauthorized)
		return
	}

	// valid groupId
	groupIdString := strings.Split(r.URL.Path, "/")[4]
	groupId, err := strconv.Atoi(groupIdString)
	if err != nil {
		http.Error(w, "groupId not valid", http.StatusBadRequest)
		return
	}
	exists, _ = rt.db.CheckIfGroupExistsByGroupId(groupId)
	if !exists {
		http.Error(w, "groupId not found", http.StatusNotFound)
		return
	}

	// userId is group partecipants
	userIdIsPartecipant, err := rt.db.CheckIfUserIsPartecipant(userId, groupId)
	if err != nil {
		http.Error(w, "Error checking userId is partecipant", http.StatusInternalServerError)
		return
	}
	if !userIdIsPartecipant {
		http.Error(w, "User Id is not a group partecipant", http.StatusUnauthorized)
		return
	}

	// load the file
	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Error loading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// create directory if not exists
	imagesDir := "images"
	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		err = os.Mkdir(imagesDir, os.ModePerm)
		if err != nil {
			http.Error(w, "Error creating images directory", http.StatusInternalServerError)
			return
		}
	}

	// create the path
	ext := filepath.Ext(header.Filename)
	if ext != ".jpeg" && ext != ".png" {
		http.Error(w, "Only .jpeg and .png image allowed", http.StatusBadRequest)
		return
	}
	fileName := fmt.Sprintf("group_%s_photo%s", groupIdString, ext)
	filePath := filepath.Join(imagesDir, fileName)

	// write the file
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
	}
	defer out.Close()

	// copy the content into the destination file
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
	}

	// update the photoUrl
	err = rt.db.UpdateGroupPhotoUrl(filePath, groupId)
	if err != nil {
		http.Error(w, "Error saving the photoUrl in the DB", http.StatusInternalServerError)
		return
	}

	// return the photoUrl
	var photoUrl PhotoUrl
	photoUrl.PhotoUrl = filePath
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(photoUrl)
	if err != nil {
		http.Error(w, "Error encoding the response", http.StatusInternalServerError)
	}
}
