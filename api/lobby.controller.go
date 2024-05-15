package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/xedom/codeduel/types"
)

func (s *Server) GetLobbyRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /lobby", convertToHandleFunc(s.handleCreateLobby))
	// router.HandleFunc("PATCH /lobby/{id}", makeHTTPHandleFunc(s.handleGetLobbyByID))
	router.HandleFunc("PATCH /lobby/{lobbyUniqueId}/submission", convertToHandleFunc(s.handleLobbyUserSubmission))
	router.HandleFunc("PATCH /lobby/{lobbyUniqueId}/endgame", convertToHandleFunc(s.handleLobbyEnd))
	router.HandleFunc("GET /lobby/results/{lobbyUniqueId}", convertToHandleFunc(s.handleGetResults))
	router.HandleFunc("OPTIONS /lobby/{lobbyUniqueId}/sharecode", convertToHandleFunc(s.handleShareCodeOptions))
	router.HandleFunc("PATCH /lobby/{lobbyUniqueId}/sharecode", convertToHandleFunc(s.handleShareCode, AuthMiddleware))
	router.HandleFunc("GET /lobby/user/{username}", convertToHandleFunc(s.handleGetMatchByUsername))

	return router
}

// @Summary		Get lobby results
// @Description	Get lobby results
// @Tags			lobby
// @Produce		json
// @Param			lobbyUniqueId	path		string	true	"Lobby unique id"
// @Success		200				{object}	types.LobbyResults
// @Failure		500				{object}	Error
// @Router			/v1/lobby/results/{lobbyUniqueId} [get]
func (s *Server) handleGetResults(w http.ResponseWriter, r *http.Request) error {
	lobbyUniqueId := r.PathValue("lobbyUniqueId")
	log.Print("[API] Getting results for lobby ", lobbyUniqueId)

	results, err := s.db.GetLobbyResults(lobbyUniqueId)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, results)
}

// @Summary		Create a new lobby
// @Description	Create a new lobby and add it to the database
// @Tags			lobby
// @Accept			json
// @Produce		json
// @Param			lobby	body	types.CreateLobbyRequest	true	"Create Lobby Request"
// @Success		204
// @Failure		500	{object}	Error
// @Router			/v1/lobby [post]
func (s *Server) handleCreateLobby(w http.ResponseWriter, r *http.Request) error {
	createLobbyPayload := &types.CreateLobbyRequest{}
	if err := json.NewDecoder(r.Body).Decode(createLobbyPayload); err != nil {
		return err
	}

	if err := s.db.CreateLobby(&types.Lobby{
		UniqueId:    createLobbyPayload.LobbyUniqueId,
		OwnerId:     createLobbyPayload.OwnerId,
		UsersId:     createLobbyPayload.UsersId,
		ChallengeId: createLobbyPayload.ChallengeId,

		Mode:             createLobbyPayload.Settings.Mode,
		MaxPlayers:       createLobbyPayload.Settings.MaxPlayers,
		GameDuration:     createLobbyPayload.Settings.GameDuration,
		AllowedLanguages: createLobbyPayload.Settings.AllowedLanguages,
	}); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// @Summary		Update lobby
// @Description	Update lobby
// @Tags			lobby
// @Produce		json
// @Param			lobby	body	types.LobbyUserSubmissionRequest	true	"Update Lobby Request"
// @Success		204
// @Failure		500	{object}	Error
// @Failure		403	{object}	Error
// @Router			/lobby/{lobbyUniqueId}/submission [patch]
func (s *Server) handleLobbyUserSubmission(w http.ResponseWriter, r *http.Request) error {
	lobbyUniqueId := r.PathValue("lobbyUniqueId")
	lobbySubmissionPayload := &types.LobbyUserSubmissionRequest{}
	if err := json.NewDecoder(r.Body).Decode(lobbySubmissionPayload); err != nil {
		return err
	}
	log.Print("[API] Lobby user submission ", lobbyUniqueId, " - ", lobbySubmissionPayload)

	lobby, err := s.db.GetLobbyByUniqueId(lobbyUniqueId)
	if err != nil {
		return err
	}

	if lobby.Ended {
		return WriteJSON(w, http.StatusForbidden, Error{Err: "The match has already ended"})
	}

	if err := s.db.UpdateLobbyUserSubmission(&types.LobbyUser{
		LobbyId:     lobby.Id,
		UserId:      lobbySubmissionPayload.UserId,
		Code:        lobbySubmissionPayload.Code,
		Language:    lobbySubmissionPayload.Language,
		TestsPassed: lobbySubmissionPayload.TestsPassed,
		SubmittedAt: lobbySubmissionPayload.Date,
	}); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

// @Summary		Update lobby
// @Description	Update lobby
// @Tags			lobby
// @Produce		json
// @Success		204
// @Failure		500	{object}	Error
// @Router			/lobby/{lobbyUniqueId}/endgame [patch]
func (s *Server) handleLobbyEnd(w http.ResponseWriter, r *http.Request) error {
	lobbyUniqueId := r.PathValue("lobbyUniqueId")
	log.Print("[API] Lobby end ", lobbyUniqueId)

	w.WriteHeader(http.StatusNoContent)
	return s.db.EndLobby(lobbyUniqueId)
}

// @Summary		Share code
// @Description	Share code
// @Tags			lobby
// @Produce		json
// @Param			lobbyUniqueId	path	string						true	"Lobby unique id"
// @Param			shareCode		body	types.ShareLobbyCodeRequest	true	"Share code request"
// @Success		204
// @Failure		500	{object}	Error
func (s *Server) handleShareCode(w http.ResponseWriter, r *http.Request) error {
	user := GetAuthUser(r)
	if user == nil {
		return WriteJSON(w, http.StatusUnauthorized, Error{Err: "Unauthorized++"})
	}

	lobbyUniqueId := r.PathValue("lobbyUniqueId")
	shareLobbyCodePayload := &types.ShareLobbyCodeRequest{}
	if err := json.NewDecoder(r.Body).Decode(shareLobbyCodePayload); err != nil {
		return err
	}

	lobby, err := s.db.GetLobbyByUniqueId(lobbyUniqueId)
	if err != nil {
		return err
	}

	log.Printf("[API] Share code for lobby %s id: %d", lobbyUniqueId, lobby.Id)

	w.WriteHeader(http.StatusNoContent)
	return s.db.UpdateShareLobbyCode(lobby.Id, user.Id, shareLobbyCodePayload.ShareCode)
}

func (s *Server) handleShareCodeOptions(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Methods", "PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// @Summary		Get match by username
// @Description	Get match by username
// @Tags			match
// @Produce		json
// @Param			username	path		string	true	"Username"
// @Success		200			{object}	[]types.SingleMatchResult
// @Failure		500			{object}	Error
// @Router			/match/user/{username} [get]
func (s *Server) handleGetMatchByUsername(w http.ResponseWriter, r *http.Request) error {
	username := r.PathValue("username")
	log.Print("[API] Fetching match for user ", username)
	matches, err := s.db.GetMatchByUsername(username)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, matches)
}
