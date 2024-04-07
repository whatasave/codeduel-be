package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xedom/codeduel/db"
	"github.com/xedom/codeduel/types"
	"github.com/xedom/codeduel/utils"
)

func (s *APIServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUsers(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateUser(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetUsers(w http.ResponseWriter, _ *http.Request) error {
	log.Print("[API] Fetching users")
	users, err := s.db.GetUsers()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, users)
}

func (s *APIServer) handleUserByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetUserByID(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteUserByID(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}

	log.Print("[API] Fetching user ", id)
	user, err := s.db.GetUserByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	createUserReq := &types.CreateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(createUserReq); err != nil {
		return err
	}

	log.Print("[API] Creating user ", createUserReq)
	user := &types.User{
		Username: createUserReq.Username,
		Email:    createUserReq.Email,
	}
	if err := s.db.CreateUser(user); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleDeleteUserByID(_ http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return err
	}

	log.Print("[API] Deleting user ", id)
	return s.db.DeleteUser(id)
}

func (s *APIServer) handleProfile(w http.ResponseWriter, r *http.Request) error {
	user, err := GetAuthUser(r, s.db)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleValidateToken(w http.ResponseWriter, r *http.Request) error {
	verifyTokenBody := &types.VerifyToken{}
	if err := json.NewDecoder(r.Body).Decode(verifyTokenBody); err != nil {
		return err
	}

	decodedUserData, err := utils.ValidateUserJWT(verifyTokenBody.JWTToken)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	// user, err := s.db.GetUserByID(userID)
	// if err != nil {
	// 	return err
	// }

	return WriteJSON(w, http.StatusOK, decodedUserData)
}

func GetAuthUser(r *http.Request, db db.DB) (*types.User, error) {
	headerUserID := r.Header.Get("x-user-id")
	userID, err := strconv.Atoi(headerUserID)
	if err != nil {
		return nil, err
	}
	
	user, err := db.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}