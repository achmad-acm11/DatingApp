package responses

import "DatingApp/entities"

type TokenResponse struct {
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   string        `json:"expires_in"`
	User        entities.User `json:"user"`
}
