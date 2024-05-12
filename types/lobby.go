package types

import "time"

type CreateLobbyRequest struct {
	LobbyId     string `json:"lobby_id"`
	OwnerId     int    `json:"owner_id"`
	UsersId     []int  `json:"users_id"`
	ChallengeId int    `json:"challenge_id"`
	Settings    struct {
		Mode             string   `json:"mode"`
		MaxPlayers       int      `json:"max_players"`
		GameDuration     int      `json:"game_duration"`
		AllowedLanguages []string `json:"allowed_languages"`
	} `json:"settings"`
}

type ShareLobbyCodeRequest struct {
	ShareCode bool `json:"share_code"`
}

type LobbyUserSubmissionRequest struct {
	UserId      int       `json:"user_id"`
	Code        string    `json:"code"`
	Language    string    `json:"language"`
	TestsPassed int       `json:"tests_passed"`
	Date        time.Time `json:"date"`
}

type Lobby struct {
	Id          int    `json:"id"`
	UniqueId    string `json:"uuid"`
	ChallengeId int    `json:"challenge_id"`
	OwnerId     int    `json:"owner_id"`
	UsersId     []int  `json:"users_id"`
	Ended       bool   `json:"ended"`

	// Settings
	Mode             string   `json:"mode"`
	MaxPlayers       int      `json:"max_players"`
	GameDuration     int      `json:"game_duration"`
	AllowedLanguages []string `json:"allowed_languages"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LobbyResults struct {
	Lobby   Lobby             `json:"lobby"`
	Results []LobbyUserResult `json:"results"`
}

type LobbyUser struct {
	Id      int `json:"id"`
	LobbyId int `json:"lobby_id"`
	UserId  int `json:"user_id"`

	Code        string `json:"code"`
	Language    string `json:"language"`
	TestsPassed int    `json:"tests_passed"`
	Rank        int    `json:"rank"`
	ShowCode    bool   `json:"show_code"`

	SubmittedAt time.Time `json:"submitted_at"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type LobbyUserResult struct {
	Id      int `json:"id"`
	LobbyId int `json:"lobby_id"`
	UserId  int `json:"user_id"`

	Code        *string `json:"code"`
	Language    *string `json:"language"`
	TestsPassed int     `json:"tests_passed"`
	ShowCode    bool    `json:"show_code"`

	SubmittedAt string `json:"submitted_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
