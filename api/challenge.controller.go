package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/xedom/codeduel/types"
)

func (s *Server) GetChallengeRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /challenge", makeHTTPHandleFunc(s.handleGetChallenges))
	router.HandleFunc("POST /challenge", makeHTTPHandleFunc(s.handleCreateChallenge))
	router.HandleFunc("GET /challenge/{id}", makeHTTPHandleFunc(s.handleGetChallengeByID))
	router.HandleFunc("PUT /challenge/{id}", makeHTTPHandleFunc(s.handleUpdateChallenge))
	router.HandleFunc("DELETE /challenge/{id}", makeHTTPHandleFunc(s.handleDeleteChallenge))
	return router
}

//	@Summary		Get all challenges
//	@Description	Get all challenges
//	@Tags			challenge
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	types.ChallengeListResponse
//	@Router			/v1/challenge [get]
func (s *Server) handleGetChallenges(w http.ResponseWriter, _ *http.Request) error {
	challenges, err := s.db.GetChallenges()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, challenges)
}

//	@Summary		Create a new challenge
//	@Description	Create a new challenge
//	@Tags			challenge
//	@Accept			json
//	@Produce		json
//	@Param			challenge	body		types.CreateChallengeRequest	true	"Create Challenge Request"
//	@Success		200			{object}	types.ChallengeResponse
//	@Router			/v1/challenge [post]
func (s *Server) handleCreateChallenge(w http.ResponseWriter, r *http.Request) error {
	createChallengeReq := &types.CreateChallengeRequest{}
	if err := json.NewDecoder(r.Body).Decode(createChallengeReq); err != nil {
		return err
	}
	
	user, err := GetUserFromDB(r, s.db)
	if err != nil {
		return err
	}

	log.Print("[API] Creating new challenge ", createChallengeReq)

	challenge := &types.Challenge{
		OwnerID: user.ID,
		Title: createChallengeReq.Title,
		Description: createChallengeReq.Description,
		Content: createChallengeReq.Content,
	}

	if err := s.db.CreateChallenge(challenge); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, challenge)
}

//	@Summary		Get challenge by ID
//	@Description	Get challenge by ID
//	@Tags			challenge
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Challenge ID"
//	@Success		200	{object}	types.Challenge
//	@Router			/v1/challenge/{id} [get]
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

//	@Summary		Update challenge by ID
//	@Description	Update challenge by ID
//	@Tags			challenge
//	@Accept			json
//	@Produce		json
//	@Param			id			path	int								true	"Challenge ID"
//	@Param			challenge	body	types.UpdateChallengeRequest	true	"Update Challenge Request"
//	@Success		200
//	@Router			/v1/challenge/{id} [put]
func (s *Server) handleUpdateChallenge(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return err
	}

	updateChallengeReq := &types.UpdateChallengeRequest{}
	if err := json.NewDecoder(r.Body).Decode(updateChallengeReq); err != nil {
		return err
	}

	log.Print("[API] Updating challenge ", id)
	challenge := &types.Challenge{
		ID: id,
		Title: updateChallengeReq.Title,
		Description: updateChallengeReq.Description,
		Content: updateChallengeReq.Content,
	}

	return s.db.UpdateChallenge(challenge)
}

//	@Summary		Delete challenge by ID
//	@Description	Delete challenge by ID
//	@Tags			challenge
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Challenge ID"
//	@Success		200
//	@Router			/v1/challenge/{id} [delete]
func (s *Server) handleDeleteChallenge(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return err
	}

	log.Print("[API] Deleting challenge ", id)
	return s.db.DeleteChallenge(id)
}
