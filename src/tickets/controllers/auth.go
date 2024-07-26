package controllers

import (
	"encoding/json"
	"strings"
	"tickets/controllers/responses"
	"tickets/utils"

	"net/http"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.StandardClaims
	Role string `json:"role,omitempty"`
}

const issuedAtLeewaySecs = 5

func (c *Claims) Valid() error {
    c.StandardClaims.IssuedAt -= issuedAtLeewaySecs
    valid := c.StandardClaims.Valid()
    c.StandardClaims.IssuedAt += issuedAtLeewaySecs
    return valid
}

func newJWKs(rawJWKS string) *keyfunc.JWKS {
	jwksJSON := json.RawMessage(rawJWKS)
	jwks, err := keyfunc.NewJSON(jwksJSON)
	if err != nil {
		panic(err)
	}
	return jwks
}

func RetrieveToken(w http.ResponseWriter, r *http.Request) *Claims {
	reqToken := r.Header.Get("Authorization")
	if len(reqToken) == 0 {
		responses.TokenIsMissing(w)
		return nil
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	tokenStr := splitToken[1]
	jwks := newJWKs(utils.Config.RawJWKS)
	tk := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, tk, jwks.Keyfunc)
	if err != nil || !token.Valid {
		responses.JwtAccessDenied(w)
		return nil
	}
	if time.Now().Unix()-tk.ExpiresAt > 0 {
		responses.TokenExpired(w)
		return nil
	}

	return tk
}

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := RetrieveToken(w, r); token != nil {
			r.Header.Set("X-User-Name", token.Subject)
			next.ServeHTTP(w, r)
		}
	})
}
