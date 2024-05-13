package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xedom/codeduel/types"
	"github.com/xedom/codeduel/utils"
)

func (s *Server) GetGithubAuthRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /auth/github", convertToHandleFunc(s.handleGithubAuth))
	router.HandleFunc("GET /auth/github/callback", convertToHandleFunc(s.handleGithubAuthCallback))
	return router
}

// @Summary		Login with GitHub
// @Description	Endpoint to log in with GitHub OAuth, it will redirect to GitHub OAuth page to authenticate
// @Tags			auth
// @Success		302
// @Router			/v1/github/auth [get]
func (s *Server) handleGithubAuth(w http.ResponseWriter, r *http.Request) error {
	urlParams := r.URL.Query()

	redirect := "https://github.com/login/oauth/authorize"

	state := genOauthState()
	http.SetCookie(w, s.createCookie("oauth_state", state, time.Now().Add(time.Minute*1)))

	// if r.FormValue("return_to") != "" {
	// 	w.Header().Set("Set-Cookie", fmt.Sprintf("return_to=%s; Path=/; HttpOnly", r.FormValue("return_to")))
	// }
	urlParams.Set("client_id", s.config.AuthGitHubClientID)
	urlParams.Set("redirect_uri", s.config.AuthGitHubClientCallbackURL)
	urlParams.Set("scope", "user:email")
	urlParams.Set("state", state) // TODO: JWT It is used to protect against cross-site request forgery attacks.
	urlParams.Set("allow_signup", "true")
	encodedParams := urlParams.Encode()

	url := fmt.Sprintf("%s?%s", redirect, encodedParams)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

// @Summary		GitHub Auth Callback
// @Description	Endpoint to handle GitHub OAuth callback, it will exchange code for access token and get user data from GitHub, then it will register a new user or login the user if it already exists. It will set a cookie with JWT token and redirect to frontend with the JWT token as a query parameter.
// @Tags			auth
// @Success		302
// @Failure		500	{object}	Error
// @Router			/v1/github/auth/callback [get]
func (s *Server) handleGithubAuthCallback(w http.ResponseWriter, r *http.Request) error {
	urlParams := r.URL.Query()
	if !urlParams.Has("code") || !urlParams.Has("state") {
		return fmt.Errorf("code or state is empty")
	}
	session_code := urlParams.Get("code")
	state := urlParams.Get("state")
	saved_state := getCookie(r, "oauth_state")

	if state != saved_state {
		return fmt.Errorf("state does not match")
	}

	githubAccessToken, err := GetGithubAccessToken(s.config.AuthGitHubClientID, s.config.AuthGitHubClientSecret, session_code, state)
	if err != nil {
		return err
	}
	log.Println("-- Github Access Token")

	githubUser, err := GetGithubUserData(githubAccessToken.AccessToken)
	if err != nil {
		return err
	}
	log.Println("-- Github User")

	if githubUser.Email == "" {
		githubEmails, err := GetGithubUserEmails(githubAccessToken.AccessToken)
		if err != nil {
			return err
		}

		// get primary email
		for _, email := range *githubEmails {
			if email.Primary {
				githubUser.Email = email.Email
				break
			}
		}
	}
	log.Println("-- Github User Email")

	// check if user exists
	auth, err := s.db.GetAuthByProviderAndID("github", fmt.Sprintf("%d", githubUser.Id))
	if err != nil {
		auth = nil
	}
	log.Println("-- Auth")

	user := &types.User{}
	var registerOrLoginError error
	if auth == nil {
		user, registerOrLoginError = RegisterGithubUser(s.db, githubUser)
	} else {
		user, registerOrLoginError = LoginGithubUser(s.db, auth)
	}
	if registerOrLoginError != nil {
		return registerOrLoginError
	}
	log.Println("-- User")

	// generating refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.Id)
	if err != nil {
		return err
	}
	log.Println("-- Refresh Token")

	err = s.db.CreateRefreshToken(user.Id, refreshToken)
	if err != nil {
		return err
	}
	log.Println("-- Refresh Token Created")

	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		return err
	}
	log.Println("-- Access Token")

	w.Header().Add("Set-Cookie", s.createCookie("refresh_token", refreshToken.Jwt, time.Unix(refreshToken.ExpiresAt, 0)).String())
	w.Header().Add("Set-Cookie", s.createCookie("access_token", accessToken.Jwt, time.Unix(accessToken.ExpiresAt, 0)).String())
	loggedInCookie := s.createCookie("logged_in", "true", time.Unix(refreshToken.ExpiresAt, 0))
	loggedInCookie.HttpOnly = false
	w.Header().Add("Set-Cookie", loggedInCookie.String())
	log.Println("-- Cookies Set")

	// redirect to frontend
	// redirectUrl := fmt.Sprintf("%s?jwt=%s", s.config.FrontendURL, refreshToken.Jwt)
	redirectUrl := s.config.FrontendURL
	if returnTo := getCookie(r, "return_to"); returnTo != "" {
		redirectUrl = fmt.Sprintf("%s/login?return_to=%s", redirectUrl, returnTo)
	}
	log.Printf("-- Redirect URL: %s\n", redirectUrl)
	w.Header().Add("Location", redirectUrl)
	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)

	return nil
}
