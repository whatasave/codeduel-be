package types

type CreateLobbyRequest struct {
	ID            int    `json:"id"`
}

type UpdateLobbyRequest struct {
	ID            int    `json:"id"`
}

type Lobby struct {
	ID            int    `json:"id"`
}
