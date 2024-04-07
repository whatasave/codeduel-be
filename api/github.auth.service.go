package api

import (
	"fmt"

	"github.com/xedom/codeduel/db"
	"github.com/xedom/codeduel/types"
	"github.com/xedom/codeduel/utils"
)

func GetGithubAccessToken(clientID, clientSecret, code, state string) (*types.GithubAccessTokenResponse, error) {
    githubUserAccessToken := &types.GithubAccessTokenResponse{}
	
	body := map[string]string{
        "client_id":     clientID,
        "client_secret": clientSecret,
        "code":          code,
        "state":         state,
    }
    
	err := utils.HttpPost("https://github.com/login/oauth/access_token", map[string]string{
		"Accept": "application/json",
		"Content-Type": "application/json",
	}, body, githubUserAccessToken)

    return githubUserAccessToken, err
}

func GetGithubUserData(accessToken string) (*types.GithubUser, error) {
    githubUserData := &types.GithubUser{}

	err := utils.HttpGet("https://api.github.com/user", map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}, githubUserData)

	return githubUserData, err
}

func GetGithubUserEmails(accessToken string) (*[]types.GithubEmails, error) {
	githubUserEmails := &[]types.GithubEmails{}

	err := utils.HttpGet("https://api.github.com/user/emails", map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", accessToken),
	}, githubUserEmails)

	return githubUserEmails, err
}

func RegisterGithubUser(db db.DB, githubUser *types.GithubUser) (*types.User, error) {
	// create user
	user := &types.User{
		Username: githubUser.Login,
		Email:    githubUser.Email,
		ImageURL: githubUser.AvatarUrl,
	}
	err := db.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// create auth
	auth := &types.AuthEntry{
		UserID:     user.ID,
		Provider:   "github",
		ProviderID: fmt.Sprintf("%d", githubUser.Id),
	}
	errAuth := db.CreateAuth(auth)
	if errAuth != nil {
		return nil, errAuth
	}

	return user, nil
}

func LoginGithubUser(db db.DB, auth *types.AuthEntry) (*types.User, error) {
	return db.GetUserByID(auth.UserID)
}