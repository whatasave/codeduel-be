package types

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type VerifyToken struct {
	JWTToken string `json:"token"`
}

type UserRequestHeader struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	// Role   string `json:"role"`
	ExpiresAt int64 `json:"expires_at"`
}

type User struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	BackgroundImg string `json:"background_img"`
	Bio           string `json:"bio"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type UserStats struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	StatsID   int    `json:"stats_id"`
	Stat      string `json:"stat"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserStatsParsed struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Stat      string `json:"stat"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Stats struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProfileResponse struct {
	*User
	Stats []*UserStatsParsed `json:"stats"`
}

type UserResponse struct {
	Name          string `json:"name"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	BackgroundImg string `json:"background_img"`
	Bio           string `json:"bio"`
	CreatedAt     string `json:"created_at"`
}
