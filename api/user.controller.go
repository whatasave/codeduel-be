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

//	@Summary		Get all users
//	@Description	Get all users from the database
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	[]types.UserResponse
//	@Failure		500	{object}	ApiError
//	@Router			/user [get]
func (s *APIServer) handleGetUsers(w http.ResponseWriter, _ *http.Request) error {
	users, err := s.db.GetUsers()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, users)
}

//	@Summary		Create a new user
//	@Description	Create a new user in the database
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		types.CreateUserRequest	true	"Create User Request"
//	@Success		200		{object}	types.User
//	@Failure		500		{object}	ApiError
//	@Router			/user [post]
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

//	@Summary		Get user by username
//	@Description	Get user by username from the database
//	@Tags			user
//	@Produce		json
//	@Param			username	path		string	true	"Username"
//	@Success		200			{object}	types.User
//	@Failure		500			{object}	ApiError
//	@Router			/user/{username} [get]
func (s *APIServer) handleGetUserByUsername(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	username := params["username"]
	log.Print("[API] Fetching user ", username)
	user, err := s.db.GetUserByUsername(username)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

//	@Summary		Delete user by username
//	@Description	Delete user by username from the database
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			username	path	string	true	"Username"
//	@Success		200
//	@Failure		500	{object}	ApiError
//	@Router			/user/{username} [delete]
func (s *APIServer) handleDeleteUserByUsername(_ http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	username := params["username"]
	log.Print("[API] Deleting user ", username)
	return s.db.DeleteUserByUsername(username)
}

//	@Summary		Get Profile
//	@Description	Get user profile when authenticated with JWT in the cookie
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	types.ProfileResponse
//	@Failure		500	{object}	ApiError
//	@Security		CookieAuth
//	@Router			/profile [get]
func (s *APIServer) handleProfile(w http.ResponseWriter, r *http.Request) error {
	user, err := GetAuthUser(r, s.db)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

//	@Summary		Validate JWT Token
//	@Description	Validate if the user JWT token is valid, and return user data. Used from other services to validate user token
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			token	body		types.VerifyToken	true	"Service token"
//	@Success		200		{object}	types.User
//	@Failure		500		{object}	ApiError
//	@Router			/validateToken [post]
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

func GetAuthUser(r *http.Request, db db.DB) (*types.ProfileResponse, error) {
	headerUserID := r.Header.Get("x-user-id")
	userID, err := strconv.Atoi(headerUserID)
	if err != nil {
		return nil, err
	}
	
	user, err := db.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	
	userStats, err := db.GetUserStats(userID)
	if err != nil {
		return nil, err
	}
	
	profile := &types.ProfileResponse{
		Stats: userStats,
		User: user,
	}
	
	return profile, nil
}
