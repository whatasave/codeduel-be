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
	router.HandleFunc("PUT /lobby/{id}", makeHTTPHandleFunc(s.handleGetLobbyByID))
	return router
}

//	@Summary		Create a new lobby
//	@Description	Create a new lobby and add it to the database
//	@Tags			lobby
//	@Accept			json
//	@Produce		json
//	@Param			lobby	body		types.CreateLobbyRequest	true	"Create Lobby Request"
//	@Success		200		{object}	types.Lobby
//	@Failure		500		{object}	Error
//	@Router			/v1/lobby [post]
func (s *Server) handleCreateLobby(w http.ResponseWriter, r *http.Request) error {
	createLobbyPayload := &types.CreateLobbyRequest{}
	if err := json.NewDecoder(r.Body).Decode(createLobbyPayload); err != nil {
		return err
	}

	log.Print("[API] Creating lobby ", createLobbyPayload)
	lobby := &types.Lobby{
		ID: 1,
	}
	// if err := s.db.CreateLobby(lobby); err != nil { return err }

	return WriteJSON(w, http.StatusOK, lobby)
}

//	@Summary		Update lobby
//	@Description	Update lobby
//	@Tags			lobby
//	@Produce		json
//	@Param			lobby	body		types.UpdateLobbyRequest	true	"Update Lobby Request"
//	@Success		200		{object}	types.Lobby
//	@Failure		500		{object}	Error
//	@Router			/v1/lobby/{id} [put]
func (s *Server) handleGetLobbyByID(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	log.Print("[API] Updating the lobby ", id)
	// lobby, err := s.db.GetLobbyByID(id)
	// if err != nil { return err }

	return WriteJSON(w, http.StatusOK, id)
}
