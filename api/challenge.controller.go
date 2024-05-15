package api

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/xedom/codeduel/types"
)

func (s *Server) GetChallengeRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /challenge", convertToHandleFunc(s.handleGetChallenges))
	router.HandleFunc("POST /challenge", convertToHandleFunc(s.handleCreateChallenge, AuthMiddleware))
	router.HandleFunc("GET /challenge/{id}", convertToHandleFunc(s.handleGetChallengeByID))
	router.HandleFunc("GET /challenge/random/full", convertToHandleFunc(s.handleGetRandomChallengeFull)) // TODO: add AuthMiddleware
	router.HandleFunc("GET /challenge/{id}/full", convertToHandleFunc(s.handleGetChallengeByIDFull, AuthMiddleware))
	router.HandleFunc("PUT /challenge/{id}", convertToHandleFunc(s.handleUpdateChallenge, AuthMiddleware))
	router.HandleFunc("DELETE /challenge/{id}", convertToHandleFunc(s.handleDeleteChallenge, AuthMiddleware))
	return router
}

// @Summary		Get all challenges
// @Description	Get all challenges
// @Tags			challenge
// @Accept			json
// @Produce		json
// @Success		200	{object}	types.ChallengeListResponse
// @Router			/v1/challenge [get]
func (s *Server) handleGetChallenges(w http.ResponseWriter, _ *http.Request) error {
	challenges, err := s.db.GetChallenges()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, challenges)
}

// @Summary		Create a new challenge
// @Description	Create a new challenge
// @Tags			challenge
// @Accept			json
// @Produce		json
// @Param			challenge	body		types.CreateChallengeRequest	true	"Create Challenge Request"
// @Success		200			{object}	types.ChallengeResponse
// @Router			/v1/challenge [post]
func (s *Server) handleCreateChallenge(w http.ResponseWriter, r *http.Request) error {
	createChallengeReq := &types.CreateChallengeRequest{}
	if err := json.NewDecoder(r.Body).Decode(createChallengeReq); err != nil {
		return err
	}

	user := GetAuthUser(r)
	if user == nil {
		return WriteJSON(w, http.StatusUnauthorized, "")
	}

	log.Print("[API] Creating new challenge ", createChallengeReq)

	challenge := &types.Challenge{
		OwnerId:     user.Id,
		Title:       createChallengeReq.Title,
		Description: createChallengeReq.Description,
		Content:     createChallengeReq.Content,
	}

	if err := s.db.CreateChallenge(challenge); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, challenge)
}

// @Summary		Get challenge by ID
// @Description	Get challenge by ID
// @Tags			challenge
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Challenge ID"
// @Success		200	{object}	types.Challenge
// @Router			/v1/challenge/{id} [get]
func (s *Server) handleGetChallengeByID(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return err
	}

	log.Print("[API] Fetching challenge ", id)
	challenge, err := s.db.GetChallengeByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, challenge)
}

// @Summary		Get challenge by ID with full details
// @Description	Get challenge by ID with full details
// @Tags			challenge
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Challenge ID"
// @Success		200	{object}	types.ChallengeFull
// @Router			/v1/challenge/{id}/full [get]
func (s *Server) handleGetChallengeByIDFull(w http.ResponseWriter, r *http.Request) error {
	// TODO check if the request is from the lobby service
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return err
	}

	log.Print("[API] Fetching challenge ", id)
	challenge, err := s.db.GetChallengeByIDFull(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, challenge)
}

// @Summary		Get random challenge with full details
// @Description	Get random challenge with full details
// @Tags			challenge
// @Accept			json
// @Produce		json
// @Success		200	{object}	types.ChallengeFull
// @Router			/v1/challenge/random/full [get]
func (s *Server) handleGetRandomChallengeFull(w http.ResponseWriter, _ *http.Request) error {
	// TODO check if the request is from the lobby service
	log.Print("[API] Fetching random challenge")

	challengesId, err := s.db.GetChallengesID()
	if err != nil {
		return err
	}

	randomId := challengesId[rand.Intn(len(challengesId))]

	challenge, err := s.db.GetChallengeByIDFull(randomId)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, challenge)
}

// @Summary		Update challenge by ID
// @Description	Update challenge by ID
// @Tags			challenge
// @Accept			json
// @Produce		json
// @Param			id			path	int								true	"Challenge ID"
// @Param			challenge	body	types.UpdateChallengeRequest	true	"Update Challenge Request"
// @Success		200
// @Router			/v1/challenge/{id} [put]
func (s *Server) handleUpdateChallenge(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return err
	}

	user := GetAuthUser(r)
	if user == nil {
		return WriteJSON(w, http.StatusUnauthorized, "")
	}
	// Unauthorized if user is not admin or owner
	if user.Role != "admin" {
		challenge, err := s.db.GetChallengeByID(id)
		if err != nil {
			return err
		}
		if challenge.OwnerId != user.Id {
			return WriteJSON(w, http.StatusUnauthorized, "")
		}
	}

	updateChallengeReq := &types.UpdateChallengeRequest{}
	if err := json.NewDecoder(r.Body).Decode(updateChallengeReq); err != nil {
		return err
	}

	log.Print("[API] Updating challenge ", id)
	challenge := &types.Challenge{
		Id:          id,
		Title:       updateChallengeReq.Title,
		Description: updateChallengeReq.Description,
		Content:     updateChallengeReq.Content,
	}

	return s.db.UpdateChallenge(challenge)
}

// @Summary		Delete challenge by ID
// @Description	Delete challenge by ID
// @Tags			challenge
// @Accept			json
// @Produce		json
// @Param			id	path	int	true	"Challenge ID"
// @Success		200
// @Router			/v1/challenge/{id} [delete]
func (s *Server) handleDeleteChallenge(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return err
	}

	user := GetAuthUser(r)
	if user == nil {
		return WriteJSON(w, http.StatusUnauthorized, "")
	}

	// Unauthorized if user is not admin or owner
	if user.Role != "admin" {
		challenge, err := s.db.GetChallengeByID(id)
		if err != nil {
			return err
		}
		if challenge.OwnerId != user.Id {
			return WriteJSON(w, http.StatusUnauthorized, "")
		}
	}

	log.Print("[API] Deleting challenge ", id)
	return s.db.DeleteChallenge(id)
}
