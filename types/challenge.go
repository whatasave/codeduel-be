package types

type Challenge struct {
	Id      int `json:"id"`
	OwnerId int `json:"owner_id"` // User.ID

	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"` // markdown maybe the link to the file

	TestCases []TestCase `json:"testCases"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type TestCase struct {
	Name   string `json:"name"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

type ChallengeFull struct {
	Id    int `json:"id"`
	Owner struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	} `json:"owner"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"` // markdown maybe the link to the file

	TestCases       []TestCase `json:"testCases"`
	HiddenTestCases []TestCase `json:"hiddenTestCases"`

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
