package utils

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/xedom/codeduel/types"
)

const jwtSecret = "yoooSuperSecret" // TODO: move to env
const expiresIn_minutes = 5

func ValidateUserJWT(tokenString string) (*types.UserRequestHeader, error) {
	token, err := ParseJWT(tokenString)
	if err != nil { return nil, err }

	// https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-token-claims#registered-claims
	claims := token.Claims.(jwt.MapClaims)
	userHeader := &types.UserRequestHeader{
		ID: int(claims["sub"].(float64)),
		Username: claims["username"].(string),
		Email: claims["email"].(string),
		ImageURL: claims["imageURL"].(string),
		// Role: claims["role"].(string),
		ExpiresAt: int64(claims["exp"].(float64)),
	}

	if err := claims.Valid(); err != nil { return nil, err }

	return userHeader, nil
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		
		return []byte(jwtSecret), nil
	})
}

type JWT struct {
	Jwt 			string `json:"jwt"`
	ExpiresAt 	int64 `json:"expires_at"`
}

func CreateJWT(user *types.User) (*JWT, error) {
	// https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-token-claims#registered-claims
	expires_at := time.Now().Add(time.Minute * expiresIn_minutes).Unix()

	claims := &jwt.MapClaims {
		"iss": "codeduel",
		"sub": user.ID,
		"exp": expires_at,

		// custom claims
		"username": user.Username,
		"email": user.Email,
		"imageURL": user.ImageURL,
		// "role": user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil { return nil, err }
	
	return &JWT{ Jwt: tokenString, ExpiresAt: expires_at, }, nil
}