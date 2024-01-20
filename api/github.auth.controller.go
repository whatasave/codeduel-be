package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/xedom/codeduel/db"
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
		urlParams := r.URL.Query()
		if !urlParams.Has("code") || !urlParams.Has("state") { return fmt.Errorf("code or state is empty") }
		code := urlParams.Get("code")
		state := urlParams.Get("state") // It is used to protect against cross-site request forgery attacks.

		client_id, err := utils.GetEnv("AUTH_GITHUB_CLIENT_ID")
		if err != nil { return err }
		client_secret, err := utils.GetEnv("AUTH_GITHUB_CLIENT_SECRET")
		if err != nil { return err }
		// client_callback_url := os.Getenv("AUTH_GITHUB_CLIENT_CALLBACK_URL")

		githubAccessToken, err := getGithubAccessToken(client_id, client_secret, code, state)
		if err != nil { return err }
		// fmt.Printf("Github Access Token: %s\n", githubAccessToken)
		
		githubUser, err := getGithubData(githubAccessToken.AccessToken)
		if err != nil { return err }
		// fmt.Printf("Github User: %+v\n", *githubUser)

		if githubUser.Email == "" {
			githubEmails, err := getGithubEmails(githubAccessToken.AccessToken)
			if err != nil { return err }
			// filter out email with primary = true

			// fmt.Printf("Github Emails: %+v\n", *githubEmails)
			for _, email := range *githubEmails {
				if email.Primary { githubUser.Email = email.Email; break }
			}
		}

		// check if user exists
		auth, err := s.db.GetAuthByProviderAndID("github", fmt.Sprintf("%d", githubUser.Id))
		if err != nil { auth = nil }

		user := &types.User{}
		var registerOrLoginError error;
		if auth == nil {
			user, registerOrLoginError = registerUser(s.db, githubUser)
		} else {
			user, registerOrLoginError = loginUser(s.db, auth)
		}
		if registerOrLoginError != nil { return registerOrLoginError }
		
		// jwt
		token, err := utils.CreateJWT(user)
		if err != nil { return err }

		// set cookie
		cookie := &http.Cookie{
			Name: "jwt",
			Value: token.Jwt,
			Domain: s.host, // TODO may cause problems
			Path: "/",
			Expires: utils.UnixTimeToTime(token.ExpiresAt),
			// MaxAge: 86400,
			HttpOnly: true,
			Secure: false,
			// SameSite: http.SameSiteStrictMode,
			// SameSite: http.SameSiteNoneMode,
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, cookie)
		// fmt.Printf("Cookie: %+v\n", cookie)
		// TODO: redirect to frontend
		// WriteJSON(w, http.StatusOK, token)
		// frontend := os.Getenv("FRONTEND_URL")
		frontendCallback := os.Getenv("FRONTEND_URL_AUTH_CALLBACK")
		redirectUrl := fmt.Sprintf("%s?jwt=%s", frontendCallback, token.Jwt)
		http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)

		return nil
	}
	
	return fmt.Errorf("method not allowed %s", r.Method)
}

func getGithubAccessToken(clientID, clientSecret, code, state string) (*types.GithubAccessTokenResponse, error) {
	requestURL := "https://github.com/login/oauth/access_token"

	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id": clientID,
		"client_secret": clientSecret,
		"code": code,
		// "redirect_uri": client_callback_url,
		"state": state,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, err := http.NewRequest("POST",requestURL,bytes.NewBuffer(requestJSON))
	if err != nil { return nil, fmt.Errorf("Request creation failed") }
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return nil, fmt.Errorf("Request failed") }

	// Convert stringified JSON to a struct object of type githubAccessTokenResponse
	respBody, _ := io.ReadAll(resp.Body)
	githubResponse := &types.GithubAccessTokenResponse{}
	json.Unmarshal(respBody, &githubResponse)

	return githubResponse, nil
}

func getGithubData(accessToken string) (*types.GithubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil { return nil, fmt.Errorf("Get Github Data API Request creation failed") }

	// Add authorization header
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return nil, fmt.Errorf("Get Github Data Request failed") }

	// Read the response as a byte slice
	respBody, _ := io.ReadAll(resp.Body)

	// Convert byte slice to string and return
	githubParsedData := &types.GithubUser{}
	jsonString := string(respBody)
	json.Unmarshal([]byte(jsonString), &githubParsedData)

	return githubParsedData, nil
}

func getGithubEmails(accessToken string) (*[]types.GithubEmails, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil { return nil, fmt.Errorf("Get Github Emails API Request creation failed") }

	// Add authorization header
	authorizationHeaderValue := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil { return nil, fmt.Errorf("Get Github Emails Request failed") }

	// Read the response as a byte slice
	resBody, _ := io.ReadAll(res.Body)

	// Convert byte slice to string and return
	githubParsedData := &[]types.GithubEmails{}
	jsonString := string(resBody)
	json.Unmarshal([]byte(jsonString), &githubParsedData)

	return githubParsedData, nil
}

func registerUser(db db.DB, githubUser *types.GithubUser) (*types.User, error) {
	// create user
	user := &types.User{
		Username: githubUser.Login,
		Email: githubUser.Email,
		ImageURL: githubUser.AvatarUrl,
	}
	err := db.CreateUser(user)
	if err != nil { return nil, err }

	// create auth
	auth := &types.AuthEntry{
		UserID: user.ID,
		Provider: "github",
		ProviderID: fmt.Sprintf("%d", githubUser.Id),
	}
	errAuth := db.CreateAuth(auth)
	if errAuth != nil { return nil, errAuth }

	return user, nil
}

func loginUser(db db.DB, auth *types.AuthEntry) (*types.User, error) {
	return db.GetUserByID(auth.UserID)
}
