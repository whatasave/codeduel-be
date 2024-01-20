package types

type AuthEntry struct {
	ID					int `json:"id"`
	UserID 			int `json:"user_id"`
	Provider 		string `json:"provider"`
	ProviderID 	string `jsong:"provider_id"`
	CreatedAt 	string `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`
}

type GithubUser struct {
	Login 							string `json:"login"`
	Id 									int `json:"id"`
	NodeId 							string `json:"node_id"`
	AvatarUrl 					string `json:"avatar_url"`
	GravatarId 					string `json:"gravatar_id"`
	Url 								string `json:"url"`
	HtmlUrl 						string `json:"html_url"`
	FollowersUrl 				string `json:"followers_url"`
	FollowingUrl 				string `json:"following_url"`
	GistsUrl 						string `json:"gists_url"`
	StarredUrl 					string `json:"starred_url"`
	SubscriptionsUrl 		string `json:"subscriptions_url"`
	Organizations_url 	string `json:"organizations_url"`
	ReposUrl 						string `json:"repos_url"`
	EventsUrl 					string `json:"events_url"`
	ReceivedEventsUrl 	string `json:"received_events_url"`
	Type 								string `json:"type"`
	SiteAdmin 					bool `json:"site_admin"`
	Name 								string `json:"name"`
	Company 						string `json:"company"`
	Blog 								string `json:"blog"`
	Location 						string `json:"location"`
	Email 							string `json:"email"`
	Hireable 						string `json:"hireable"`
	Bio 								string `json:"bio"`
	TwitterUsername 		string `json:"twitter_username"`
	PublicRepos 				int `json:"public_repos"`
	PublicGists 				int `json:"public_gists"`
	Followers 					int `json:"followers"`
	Following 					int `json:"following"`
	CreatedAt 					string `json:"created_at"`
	UpdatedAt 					string `json:"updated_at"`
}

type GithubEmails struct {
	Email 				string `json:"email"`
	Verified 			bool `json:"verified"`
	Primary 			bool `json:"primary"`
	Visibility 		string `json:"visibility"`
}


type GithubAccessTokenResponse struct {
	AccessToken 	string `json:"access_token"`
	TokenType   	string `json:"token_type"`
	Scope       	string `json:"scope"`
}
