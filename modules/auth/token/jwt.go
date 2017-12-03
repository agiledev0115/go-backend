package token

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claim struct {
	Issuer   string
	Audience string
	Subject  string
}

type AccessToken struct {
	Token string `json:"access_token"`
}

type AccessTokenResponse struct {
	Error       error
	AccessToken *AccessToken
}

func GenerateAccessToken(signKey *rsa.PrivateKey, cl Claim, tokenAge time.Duration) <-chan AccessTokenResponse {
	result := make(chan AccessTokenResponse)
	go func() {
		defer close(result)
		token := jwt.New(jwt.SigningMethodRS256)
		claims := make(jwt.MapClaims)
		claims["iss"] = cl.Issuer
		claims["aud"] = cl.Audience
		claims["exp"] = time.Now().Add(tokenAge).Unix()
		claims["iat"] = time.Now().Unix()
		claims["sub"] = cl.Subject
		token.Claims = claims

		tokenString, err := token.SignedString(signKey)
		if err != nil {
			result <- AccessTokenResponse{Error: err, AccessToken: nil}
			return
		}
		result <- AccessTokenResponse{Error: nil, AccessToken: &AccessToken{Token: fmt.Sprintf("Bearer %v", tokenString)}}
	}()
	return result
}