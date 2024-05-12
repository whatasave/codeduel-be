package api

import (
	"fmt"
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
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  state,
		Path:   "/",
		MaxAge: 60,
	})

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
	returnTo := getCookie(r, "return_to")

	if state != saved_state {
		return fmt.Errorf("state does not match")
	}

	githubAccessToken, err := GetGithubAccessToken(s.config.AuthGitHubClientID, s.config.AuthGitHubClientSecret, session_code, state)
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

	// generating refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.Id)
	if err != nil {
		return err
	}

	err = s.db.CreateRefreshToken(user.Id, refreshToken)
	if err != nil {
		return err
	}

	// set cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "refresh_token",
		Value:   refreshToken.Jwt,
		Domain:  s.config.CookieDomain,
		Path:    s.config.CookiePath,
		Expires: time.Unix(refreshToken.ExpiresAt, 0),
		// MaxAge: s.config.CookieMaxAge,
		HttpOnly: s.config.CookieHTTPOnly,
		Secure:   s.config.CookieSecure,
		// SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSiteNoneMode,
		SameSite: http.SameSiteLaxMode,
	})

	// redirect to frontend
	// redirectUrl := fmt.Sprintf("%s?jwt=%s", s.config.FrontendURL, refreshToken.Jwt)
	redirectUrl := s.config.FrontendURL
	if returnTo != "" {
		redirectUrl += fmt.Sprintf("&return_to=%s", returnTo)
	}
	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)

	return nil
}
