package types

type Challenge struct {
	Id          int    `json:"id"`
	OwnerId     int    `json:"owner_id"` // User.ID
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"` // markdown maybe the link to the file

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateChallengeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type UpdateChallengeRequest struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

type ChallengeResponse struct {
	Id          int    `json:"id"`
	OwnerId     int    `json:"owner_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ChallengeListResponse struct {
	Challenges []ChallengeResponse `json:"challenges"`
}
