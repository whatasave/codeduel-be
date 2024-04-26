package types

import "time"

type CreateLobbyRequest struct {
	LobbyId     string `json:"lobby_id"`
	OwnerId     int    `json:"owner_id"`
	UsersId     []int  `json:"users_id"`
	ChallengeId int    `json:"challenge_id"`
	Settings    struct {
		MaxPlayers       int           `json:"max_players"`
		GameDuration     time.Duration `json:"game_duration"`
		AllowedLanguages []string      `json:"allowed_languages"`
	} `json:"settings"`
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
	Status      string `json:"status"`

	// Settings
	MaxPlayers       int           `json:"max_players"`
	GameDuration     time.Duration `json:"game_duration"`
	AllowedLanguages []string      `json:"allowed_languages"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LobbyUser struct {
	Id      int `json:"id"`
	LobbyId int `json:"lobby_id"`
	UserId  int `json:"user_id"`

	Code           string    `json:"code"`
	Language       string    `json:"language"`
	TestsPassed    int       `json:"tests_passed"`
	SubmissionDate time.Time `json:"submission_date"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
