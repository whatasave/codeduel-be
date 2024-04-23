package api

import (
	"fmt"
	"net/http"

	"github.com/xedom/codeduel/types"
	"github.com/xedom/codeduel/utils"
)

func (s *Server) GetGithubAuthRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /auth/github", makeHTTPHandleFunc(s.handleGithubAuth))
	router.HandleFunc("GET /auth/github/callback", makeHTTPHandleFunc(s.handleGithubAuthCallback))
	return router
}

//	@Summary		Login with GitHub
//	@Description	Endpoint to log in with GitHub OAuth, it will redirect to GitHub OAuth page to authenticate
//	@Tags			auth
//	@Success		302
//	@Router			/v1/github/auth [get]
func (s *Server) handleGithubAuth(w http.ResponseWriter, r *http.Request) error {
	urlParams := r.URL.Query()

	redirect := "https://github.com/login/oauth/authorize"

	urlParams.Add("client_id", s.config.AuthGitHubClientID)
	urlParams.Add("redirect_uri", s.config.AuthGitHubClientCallbackURL)
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

//	@Summary		GitHub Auth Callback
//	@Description	Endpoint to handle GitHub OAuth callback, it will exchange code for access token and get user data from GitHub, then it will register a new user or login the user if it already exists. It will set a cookie with JWT token and redirect to frontend with the JWT token as a query parameter.
//	@Tags			auth
//	@Success		302
//	@Failure		500	{object} Error
//	@Router			/v1/github/auth/callback [get]
func (s *Server) handleGithubAuthCallback(w http.ResponseWriter, r *http.Request) error {
	urlParams := r.URL.Query()
	if !urlParams.Has("code") || !urlParams.Has("state") {
		return fmt.Errorf("code or state is empty")
	}
	code := urlParams.Get("code")
	state := urlParams.Get("state") // It is used to protect against cross-site request forgery attacks.

	githubAccessToken, err := GetGithubAccessToken(s.config.AuthGitHubClientID, s.config.AuthGitHubClientSecret, code, state)
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
		Domain:  s.config.CookieDomain,
		Path:    s.config.CookiePath,
		Expires: utils.UnixTimeToTime(token.ExpiresAt),
		// MaxAge: 86400,
		HttpOnly: s.config.CookieHTTPOnly, //true,
		Secure:   s.config.CookieSecure,   //false,
		// SameSite: http.SameSiteStrictMode,
		// SameSite: http.SameSiteNoneMode,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
	// fmt.Printf("Cookie: %+v\n", cookie)
	// TODO: redirect to frontend
	// WriteJSON(w, http.StatusOK, token)
	redirectUrl := fmt.Sprintf("%s?jwt=%s", s.config.FrontendURLAuthCallback, token.Jwt)
	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)

	return nil
}
