package utils

import (
	"fmt"
	"reflect"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/xedom/codeduel/types"
)

var (
	jwtSecret                    string
	expiresInMinutes             int
	refreshTokenExpiresInMinutes int
)

func init() {
	config = LoadConfig()

	jwtSecret = config.JWTSecret
	expiresInMinutes = config.JWTExpiresInMinutes
	refreshTokenExpiresInMinutes = config.JWTRefreshTokenExpiresInMinutes
}

func ValidateUserJWT(tokenString string) (*types.UserRequestHeader, error) {
	tokenClaims, err := ParseJWT(tokenString)
	if err != nil {
		return nil, err
	}

	// https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-token-claims#registered-claims
	if err := tokenClaims.Valid(); err != nil {
		return nil, err
	}

	return &types.UserRequestHeader{
		Id:        int((*tokenClaims)["sub"].(float64)),
		Username:  (*tokenClaims)["username"].(string),
		Email:     (*tokenClaims)["email"].(string),
		Avatar:    (*tokenClaims)["avatar"].(string),
		Role:      (*tokenClaims)["role"].(string),
		ExpiresAt: int64((*tokenClaims)["exp"].(float64)),
	}, nil
}

func ValidateAndParseJWT(tokenString string, structToParse interface{}) error {
	reflectValue := reflect.ValueOf(structToParse)
	if reflectValue.Kind() != reflect.Ptr {
		return fmt.Errorf("structToParse must be a pointer")
	}

	tokenClaims, err := ParseJWT(tokenString)
	if err != nil {
		return err
	}

	// https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-token-claims#registered-claims
	if err := tokenClaims.Valid(); err != nil {
		return err
	}

	// reflectValue.Elem().Set(reflect.ValueOf(*tokenClaims))
	// using json tags to parse the jwt claims
	if err := MapClaimsToStruct(*tokenClaims, structToParse); err != nil {
		return err
	}

	return nil
}

func MapClaimsToStruct(claims jwt.MapClaims, structToParse interface{}) error {
	reflectValue := reflect.ValueOf(structToParse)
	if reflectValue.Kind() != reflect.Ptr {
		return fmt.Errorf("structToParse must be a pointer")
	}

	reflectValue = reflectValue.Elem()
	reflectType := reflectValue.Type()

	for i := 0; i < reflectValue.NumField(); i++ {
		fieldValue := reflectValue.Field(i)
		fieldType := reflectType.Field(i)

		// check if the field is exported
		if fieldType.PkgPath != "" {
			continue
		}

		// check if the field is in the claims
		claim, ok := claims[fieldType.Tag.Get("jwt")]
		if !ok {
			continue
		}

		// cast the claim to the field type
		switch fieldType.Type.Kind() {
		case reflect.Int:
			claim = int(claim.(float64))
		case reflect.Int64:
			claim = int64(claim.(float64))
		case reflect.String:
			claim = claim.(string)
		}

		// check if the field is a pointer
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue.Set(reflect.ValueOf(&claim))
		} else {
			fieldValue.Set(reflect.ValueOf(claim))
		}
	}

	return nil

}

func ParseJWT(tokenString string) (*jwt.MapClaims, error) {
	jwtParsed, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims := jwtParsed.Claims.(jwt.MapClaims)

	return &claims, nil

}

type JWT struct {
	Jwt       string `json:"jwt"`
	ExpiresAt int64  `json:"expires_at"`
}

func CreateJWT(claims *jwt.MapClaims) (*JWT, error) {
	// https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-token-claims#registered-claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	exp := (*claims)["exp"].(int64)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	return &JWT{Jwt: tokenString, ExpiresAt: exp}, nil
}

func GenerateRefreshToken(userId int) (*JWT, error) {
	return CreateJWT(&jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Minute * time.Duration(refreshTokenExpiresInMinutes)).Unix(),
	})
}

func GenerateAccessToken(user *types.User) (*JWT, error) {
	return CreateJWT(&jwt.MapClaims{
		"iss": "codeduel",
		"sub": user.Id,
		"exp": time.Now().Add(time.Minute * time.Duration(expiresInMinutes)).Unix(),

		// custom claims
		"username": user.Username,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"role":     user.Role,
	})
}
