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
	router.HandleFunc("POST /lobby/{lobbyId}/submission", makeHTTPHandleFunc(s.handleLobbyUserSubmission))
	router.HandleFunc("PATCH /lobby/{lobbyId}/endgame", makeHTTPHandleFunc(s.handleLobbyEnd))
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
	lobby := &types.Lobby{Id: 1}

	return WriteJSON(w, http.StatusOK, lobby)
}

// @Summary		Update lobby
// @Description	Update lobby
// @Tags			lobby
// @Produce		json
// @Param			lobby	body	types.UpdateLobbyRequest	true	"Update Lobby Request"
// @Success		204
// @Failure		500	{object}	Error
// @Router			/lobby/{lobbyId}/submission [patch]
func (s *Server) handleLobbyUserSubmission(w http.ResponseWriter, r *http.Request) error {
	lobbyId := r.PathValue("lobbyId")
	log.Print("[API] Lobby user submission ", lobbyId)

	return WriteJSON(w, http.StatusNoContent, lobbyId)
}

// @Summary		Update lobby
// @Description	Update lobby
// @Tags			lobby
// @Produce		json
// @Success		204
// @Failure		500	{object}	Error
// @Router			/lobby/{lobbyId}/endgame [patch]
func (s *Server) handleLobbyEnd(w http.ResponseWriter, r *http.Request) error {
	lobbyId := r.PathValue("lobbyId")
	log.Print("[API] Lobby end ", lobbyId)

	return WriteJSON(w, http.StatusNoContent, lobbyId)
}
