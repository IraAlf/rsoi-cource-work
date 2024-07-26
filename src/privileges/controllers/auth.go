package controllers

import (
	"encoding/json"
	"log"
	"privileges/controllers/responses"
	"strings"

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

var jwtKey = []byte("your-256-bit-secret")

func (c *Claims) Valid() error {
	c.StandardClaims.IssuedAt -= issuedAtLeewaySecs
	valid := c.StandardClaims.Valid()
	c.StandardClaims.IssuedAt += issuedAtLeewaySecs
	return valid
}

func newJWKs(rawJWKS string) *keyfunc.JWKS {
	jwksJSON := json.RawMessage(rawJWKS)
	jwks, err := keyfunc.NewJSON(jwksJSON)
	log.Printf("Privileges: token: %s ", jwks)
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		panic(err)
	}
	return jwks
}

func RetrieveToken(w http.ResponseWriter, r *http.Request) *Claims {
	reqToken := r.Header.Get("Authorization")
	log.Printf("Privileges: token: %s ", reqToken)
	if len(reqToken) == 0 {
		responses.TokenIsMissing(w)
		return nil
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	tokenStr := splitToken[1]

	//jwks := newJWKs(utils.Config.RawJWKS)
	//jwks := newJWKs(`{"keys":[{"kty":"oct","kid":"your-key-id","k":"MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQClPpTCCMRhCTWDExdmXXH+AZVNsIX4VrbI0jvUZmfSEZNNvpyQ48SeA2xF3hL3iEjlXIqa0lCs7wxn+Rk11Ezi82yLRubK+/emP1JfsCrx0WnZEoUU0SwgIEE9Igb1jMBHZvTYPmNDz/B2ZnmXQ481gSWKvsydI2JJYEj14bNrRwIDAQAB"}]}`)
	tk := &Claims{}

	// token, err := jwt.ParseWithClaims(tokenStr, tk, jwks.Keyfunc)

	token, err := jwt.ParseWithClaims(tokenStr, tk, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtKey, nil
	})

	// Log the JWT header to check for the presence of 'kid'

	if err != nil || !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				log.Printf("Malformed token: %v", err)
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				log.Printf("Token expired or not valid yet: %v", err)
			} else {
				log.Printf("Invalid token: %v", err)
			}
		} else {
			log.Printf("Error parsing token: %v", err)
		}
		log.Printf("JwtAccessDenied", err)
		responses.JwtAccessDenied(w)
		return nil
	}
	if time.Now().Unix()-tk.ExpiresAt > 0 {
		log.Printf("ExpiresAt")
		responses.TokenExpired(w)
		return nil
	}
	log.Printf("Privileges: token: %s ", tk)
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
