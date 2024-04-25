package types

type Match struct {
	Id          int   `json:"id"`
	OwnerId     int   `json:"owner_id"`     // User.ID
	ChallengeId int   `json:"challenge_id"` // Challenge.ID
	ModeId      int   `json:"mode_id"`      // Mode.ID
	MaxUsers    int   `json:"max_users"`
	MaxDuration int   `json:"max_time"`
	AllowedLang []int `json:"allowed_lang"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MatchUserLink struct {
	Id            int    `json:"id"`
	MatchId       int    `json:"match_id"`     // Match.ID
	UserId        int    `json:"user_id"`      // User.ID
	StatusId      int    `json:"status"`       // Status.ID
	MatchStatusId int    `json:"match_status"` // MatchStatus.ID
	Code          string `json:"code"`
	LanguageId    int    `json:"language_id"` // Language.ID
	Rank          int    `json:"rank"`
	Duration      int    `json:"duration"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Mode struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Status 0: not ready, 1: ready, 2: in match, 3: finished
type Status struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// MatchStatus 0: starting, 1: ongoing, 2: finished
type MatchStatus struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Language struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
