package types

type Match struct {
	ID 					int `json:"id"`
	OwnerID 		int `json:"owner_id"` // User.ID
	ChallengeID int `json:"challenge_id"` // Challenge.ID
	ModeID 			int `json:"mode_id"` // Mode.ID
	MaxUsers 		int `json:"max_users"`
	MaxDuration	int `json:"max_time"`
	AllowedLang []int `json:"allowed_lang"`
	CreatedAt 	string `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`
}

type MatchUserLink struct {
	ID 						int `json:"id"`
	MatchID 			int `json:"match_id"` // Match.ID
	UserID 				int `json:"user_id"` // User.ID
	StatusID 			int `json:"status"`  // Status.ID
	MatchStatusID int `json:"match_status"` // MatchStatus.ID
	Code 					string `json:"code"`
	LanguageID 		int `json:"language_id"` // Language.ID
	Rank 					int `json:"rank"`
	Duration 			int `json:"duration"`
	CreatedAt 		string `json:"created_at"`
	UpdatedAt 		string `json:"updated_at"`
}

type Challenge struct {
	ID 					int `json:"id"`
	OwnerID 		int `json:"owner_id"` // User.ID
	Title 			string `json:"title"`
	Description string `json:"description"`
	Content 		string `json:"content"` // markdown maybe the link to the file
	CreatedAt 	string `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`
}

type Mode struct {
	ID 					int `json:"id"`
	Name 				string `json:"name"`
	Description string `json:"description"`
}

// 0: not ready, 1: ready, 2: in match, 3: finished
type Status struct {
	ID 					int `json:"id"`
	Name 				string `json:"name"`
}

// 0: starting, 1: ongoing, 2: finished
type MatchStatus struct {
	ID 					int `json:"id"`
	Name 				string `json:"name"`
}

type Language struct {
	ID 					int `json:"id"`
	Name 				string `json:"name"`
}