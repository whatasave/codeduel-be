package types

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type VerifyToken struct {
	JWTToken string `json:"token"`
}

type UserRequestHeader struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expires_at"`
}

type User struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	BackgroundImg string `json:"background_img"`
	Bio           string `json:"bio"`
	Role          string `json:"role"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserStats struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id"`
	StatsId int    `json:"stats_id"`
	Stat    string `json:"stat"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserStatsParsed struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Stat string `json:"stat"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Stats struct {
	Id   int    `json:"id"`
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
	Role          string `json:"role"`
	CreatedAt     string `json:"created_at"`
}

type RefreshTokenPayload struct {
	UserID    int   `json:"user_id" jwt:"sub"`
	ExpiresAt int64 `json:"expires_at" jwt:"exp"`
}
