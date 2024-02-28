package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/xedom/codeduel/types"
	"github.com/xedom/codeduel/utils"
)

func (s *APIServer) handleGithubAuth(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		// redirect to github auth
		urlParams := r.URL.Query()

		client_id := os.Getenv("AUTH_GITHUB_CLIENT_ID")
		// client_secret := os.Getenv("AUTH_GITHUB_CLIENT_SECRET")
		client_callback_url := os.Getenv("AUTH_GITHUB_CLIENT_CALLBACK_URL")
		redirect := "https://github.com/login/oauth/authorize"

		urlParams.Add("client_id", client_id)
		urlParams.Add("redirect_uri", client_callback_url)
		urlParams.Add("return_to", "/frontend")
		urlParams.Add("response_type", "code")
		urlParams.Add("scope", "user:email")
		urlParams.Add("state", "an_unguessable_random_string") // TODO: JWT It is used to protect against cross-site request forgery attacks.
		urlParams.Add("allow_signup", "true")
		encodedParams := urlParams.Encode()

		url := fmt.Sprintf("%s?%s", redirect, encodedParams)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGithubAuthCallback(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGithubAuthGetRequest(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGithubAuthGetRequest(w http.ResponseWriter, r *http.Request) error {

	urlParams := r.URL.Query()
	if !urlParams.Has("code") || !urlParams.Has("state") {
		return fmt.Errorf("code or state is empty")
	}
	code := urlParams.Get("code")
	state := urlParams.Get("state") // It is used to protect against cross-site request forgery attacks.

	client_id, err := utils.GetEnv("AUTH_GITHUB_CLIENT_ID")
	if err != nil {
		return err
	}
	client_secret, err := utils.GetEnv("AUTH_GITHUB_CLIENT_SECRET")
	if err != nil {
		return err
	}
	// client_callback_url := os.Getenv("AUTH_GITHUB_CLIENT_CALLBACK_URL")
	frontendCallback, err := utils.GetEnv("FRONTEND_URL_AUTH_CALLBACK")
	if err != nil {
		return err
	}

	githubAccessToken, err := GetGithubAccessToken(client_id, client_secret, code, state)
	if err != nil {
		return err
	}
	// fmt.Printf("Github Access Token: %s\n", githubAccessToken)

	githubUser, err := GetGithubUserData(githubAccessToken.AccessToken)
	if err != nil {
		return err
	}
	// fmt.Printf("Github User: %+v\n", *githubUser)

	if githubUser.Email == "" {
		githubEmails, err := GetGithubUserEmails(githubAccessToken.AccessToken)
		if err != nil {
			return err
		}

		// get primary email
		// fmt.Printf("Github Emails: %+v\n", *githubEmails)
		for _, email := range *githubEmails {
			if email.Primary {
				githubUser.Email = email.Email
				break
			}
		}
	}

	// check if user exists
	auth, err := s.db.GetAuthByProviderAndID("github", fmt.Sprintf("%d", githubUser.Id))
	if err != nil {
		auth = nil
	}

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

	// generate jwt
	token, err := utils.CreateJWT(user)
	if err != nil {
		return err
	}

	// set cookie
	cookie := &http.Cookie{
		Name:    "jwt",
		Value:   token.Jwt,
		Domain:  s.host, // TODO may cause problems
		Path:    "/",
		Expires: utils.UnixTimeToTime(token.ExpiresAt),
		// MaxAge: 86400,
		HttpOnly: true,
		Secure:   false,
		// SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSiteNoneMode,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	// fmt.Printf("Cookie: %+v\n", cookie)
	// TODO: redirect to frontend
	// WriteJSON(w, http.StatusOK, token)
	// frontend := os.Getenv("FRONTEND_URL")
	redirectUrl := fmt.Sprintf("%s?jwt=%s", frontendCallback, token.Jwt)
	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)

	return nil
}

