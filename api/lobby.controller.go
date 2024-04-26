package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/xedom/codeduel/types"
)

func (s *Server) GetLobbyRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /lobby", makeHTTPHandleFunc(s.handleCreateLobby))
	// router.HandleFunc("PATCH /lobby/{id}", makeHTTPHandleFunc(s.handleGetLobbyByID))
	router.HandleFunc("POST /lobby/{lobbyUniqueId}/submission", makeHTTPHandleFunc(s.handleLobbyUserSubmission))
	router.HandleFunc("PATCH /lobby/{lobbyUniqueId}/endgame", makeHTTPHandleFunc(s.handleLobbyEnd))
	return router
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
	log.Print("[API] Creating lobby ", createLobbyPayload)

	if err := s.db.CreateLobby(&types.Lobby{
		UniqueId:    createLobbyPayload.LobbyId,
		OwnerId:     createLobbyPayload.OwnerId,
		UsersId:     createLobbyPayload.UsersId,
		ChallengeId: createLobbyPayload.ChallengeId,

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
// @Router			/lobby/{lobbyUniqueId}/submission [post]
func (s *Server) handleLobbyUserSubmission(w http.ResponseWriter, r *http.Request) error {
	lobbyUniqueId := r.PathValue("lobbyUniqueId")
	lobbySubmissionPayload := &types.LobbyUserSubmissionRequest{}
	if err := json.NewDecoder(r.Body).Decode(lobbySubmissionPayload); err != nil {
		return err
	}
	log.Print("[API] Lobby user submission ", lobbyUniqueId, lobbySubmissionPayload)

	lobby, err := s.db.GetLobbyByUniqueId(lobbyUniqueId)
	if err != nil {
		return err
	}

	if lobby.Status != "open" {
		return WriteJSON(w, http.StatusForbidden, Error{Err: "Lobby is not open"})
	}

	if err := s.db.CreateLobbyUserSubmission(&types.LobbyUser{
		LobbyId:        lobby.Id,
		UserId:         lobbySubmissionPayload.UserId,
		Code:           lobbySubmissionPayload.Code,
		Language:       lobbySubmissionPayload.Language,
		TestsPassed:    lobbySubmissionPayload.TestsPassed,
		SubmissionDate: lobbySubmissionPayload.Date,
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
