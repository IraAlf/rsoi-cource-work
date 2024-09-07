package controllers

import (
	"encoding/json"
	"fmt"
	"identity-provider/controllers/responses"
	"identity-provider/models"
	"identity-provider/objects"
	"log"
	"net/http"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type auhtCtrl struct {
	user *models.UserModel
}

var jwtKey = []byte("your-256-bit-secret")

type Claims struct {
	Subject string `json:"sub"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func InitAuth(r *mux.Router, auth *models.UserModel) {
	ctrl := &auhtCtrl{auth}
	r.HandleFunc("/register", ctrl.register).Methods("POST")
	r.HandleFunc("/authorize", ctrl.authorize).Methods("POST")
	r.HandleFunc("/manage/health", ctrl.GetHealth).Methods("GET")
}

func (h *auhtCtrl) GetHealth(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}

type Token struct {
	jwt.StandardClaims
	Role string `json:"role,omitempty"`
}

const issuedAtLeewaySecs = 5

func (c *Token) Valid() error {
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

func (ctrl *auhtCtrl) register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Identity-provider: register: %s", r.Body)
	requestBody := new(objects.UserCreateRequest)
	json.NewDecoder(r.Body).Decode(requestBody)
	log.Printf("Identity-provider: register: %s", requestBody)
	createdUser, err := ctrl.user.RegisterUser(requestBody.Profile.Firstname,
		requestBody.Profile.Lastname,
		requestBody.Profile.Login,
		requestBody.Credentials.Password.Value,
		requestBody.Profile.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// expirationTime := time.Now().Add(60 * time.Minute)
	// claims := &Claims{
	// 	Subject: requestBody.Profile.Login,
	// 	Role:    requestBody.Profile.UserType,
	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		// In JWT, the expiry time is expressed as unix milliseconds
	// 		ExpiresAt: jwt.NewNumericDate(expirationTime),
	// 	},
	// }

	// // Declare the token with the algorithm used for signing, and the claims
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// // Create the JWT string
	// tokenString, err := token.SignedString(jwtKey)

	// w.WriteHeader("Authorization", tokenString)

	responses.JsonSuccess(w, createdUser)
}

func (ctrl *auhtCtrl) authorize(w http.ResponseWriter, r *http.Request) {
	var credentials objects.AuthRequest
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userFromDb, err := ctrl.user.GetUser(credentials.Login)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	log.Printf("Identity-provider: register: %s %s", userFromDb.Password, credentials.Password)

	// Get the expected password from our in memory map
	expectedPassword := userFromDb.Password

	// If a password exists for the given user
	// AND, if it is the same as the password we received, then we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if expectedPassword != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(60 * time.Minute)
	// Create the JWT claims, which includes the login and expiry time
	claims := &Claims{
		Subject: credentials.Login,
		Role:    userFromDb.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	fmt.Println("userFromDb.UserType %s", userFromDb.UserType)

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Add 'kid' to the token header

	//token.Header["kid"] = "your-key-id"

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Identity-provider: error: %s", err)
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := &objects.AuthResponse{
		ExpiresIn:   int(expirationTime.Unix()),
		AccessToken: tokenString,
		Role:        userFromDb.UserType,
	}

	log.Printf("Identity-provider: token: %s ", response.AccessToken)

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	responses.JsonSuccess(w, response)
}
