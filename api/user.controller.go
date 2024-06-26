package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/xedom/codeduel/types"
)

func (s *Server) GetUserRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /user", convertToHandleFunc(s.handleGetUsers))
	router.HandleFunc("POST /user", convertToHandleFunc(s.handleCreateUser, AuthMiddleware))
	router.HandleFunc("GET /user/{username}", convertToHandleFunc(s.handleGetUserByUsername))
	router.HandleFunc("DELETE /user/{username}", convertToHandleFunc(s.handleDeleteUserByUsername, AuthMiddleware))
	router.HandleFunc("GET /user/profile", convertToHandleFunc(s.handleProfile, AuthMiddleware))

	return router
}

// @Summary		Get all users
// @Description	Get all users from the database
// @Tags			user
// @Produce		json
// @Success		200	{object}	[]types.UserResponse
// @Failure		500	{object}	Error
// @Router			/v1/user [get]
func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) error {
	users := []*types.UserResponse{}

	if r.URL.Query().Get("id") != "" {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			return err
		}
		user, err := s.db.GetUserByID(id)
		if err != nil {
			return err
		}
		users = append(users, &types.UserResponse{
			Name:          user.Name,
			Username:      user.Username,
			Avatar:        user.Avatar,
			BackgroundImg: user.BackgroundImg,
			Bio:           user.Bio,
			Role:          user.Role,
			CreatedAt:     user.CreatedAt,
		})
	} else {
		dbUsers, err := s.db.GetUsers()
		if err != nil {
			return err
		}
		users = dbUsers
	}

	return WriteJSON(w, http.StatusOK, users)
}

// @Summary		Create a new user
// @Description	Create a new user in the database
// @Tags			user
// @Accept			json
// @Produce		json
// @Param			user	body		types.CreateUserRequest	true	"Create User Request"
// @Success		200		{object}	types.User
// @Failure		500		{object}	Error
// @Router			/v1/user [post]
func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
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

// @Summary		Get user by username
// @Description	Get user by username from the database
// @Tags			user
// @Produce		json
// @Param			username	path		string	true	"Username"
// @Success		200			{object}	types.User
// @Failure		500			{object}	Error
// @Router			/v1/user/{username} [get]
func (s *Server) handleGetUserByUsername(w http.ResponseWriter, r *http.Request) error {
	username := r.PathValue("username")
	log.Print("[API] Fetching user ", username)
	user, err := s.db.GetUserByUsername(username)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

// @Summary		Delete user by username
// @Description	Delete user by username from the database
// @Tags			user
// @Accept			json
// @Produce		json
// @Param			username	path	string	true	"Username"
// @Success		200
// @Failure		500	{object}	Error
// @Router			/v1/user/{username} [delete]
func (s *Server) handleDeleteUserByUsername(_ http.ResponseWriter, r *http.Request) error {
	username := r.PathValue("username")
	log.Print("[API] Deleting user ", username)
	return s.db.DeleteUserByUsername(username)
}

// @Summary		Get Profile
// @Description	Get user profile when authenticated with JWT in the cookie
// @Tags			user
// @Produce		json
// @Success		200	{object}	types.UserRequestHeader
// @Failure		500	{object}	Error
// @Security		CookieAuth
// @Router			/v1/user/profile [get]
func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) error {
	user := GetAuthUser(r)

	return WriteJSON(w, http.StatusOK, user)
}
