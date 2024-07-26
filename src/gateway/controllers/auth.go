package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gateway/controllers/responses"
	"gateway/models"
	"gateway/objects"
	"gateway/utils"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
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

func RetrieveToken(w http.ResponseWriter, r *http.Request) (*Claims, error) {
	reqToken := r.Header.Get("Authorization")
	if len(reqToken) == 0 {
		responses.TokenIsMissing(w)
		return nil, fmt.Errorf("TokenIsMissing")
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	tokenStr := splitToken[1]
	jwks := newJWKs(utils.Config.RawJWKS)
	tk := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, tk, jwks.Keyfunc)
	if err != nil || !token.Valid {
		responses.JwtAccessDenied(w)
		return nil, fmt.Errorf("JwtAccessDenied")
	}
	if time.Now().Unix()-tk.ExpiresAt > 0 {
		responses.TokenExpired(w)
		return nil, fmt.Errorf("TokenExpired")
	}

	return tk, nil
}

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := RetrieveToken(w, r)
		if err != nil {
			log.Printf("Token error: %s", err.Error())
			return
		}

		r.Header.Set("X-User-Name", token.Subject)
		r.Header.Set("X-User-Role", token.Role)
		next.ServeHTTP(w, r)
	})
}

type authCtrl struct {
	client     *http.Client
	privileges *models.PrivilegesM
}

func InitAuth(r *mux.Router, client *http.Client, privileges *models.PrivilegesM) {
	ctrl := &authCtrl{client, privileges}
	r.HandleFunc("/register", ctrl.register).Methods("POST")
	r.HandleFunc("/authorize", ctrl.authorize).Methods("POST")
}

func (ctrl *authCtrl) register(w http.ResponseWriter, r *http.Request) {
	token, err := RetrieveToken(w, r)
	if err != nil {
		log.Printf("failed to RetrieveToken: %s", err.Error())
		return
	}

	if token.Role != "admin" {
		responses.ForbiddenMsg(w, fmt.Sprintf("not allowed for %s role", token.Role))
		return
	}

	req_body := new(objects.UserCreateRequest)
	err = json.NewDecoder(r.Body).Decode(req_body)
	log.Printf("creating new account: %v", req_body)
	if err != nil {
		log.Printf("failed to parse body: %s", err.Error())
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("sakura response: %q", r.Body)

		responses.ValidationErrorResponse(w, err.Error())
		return
	}

	req, shouldReturn := ctrl.makeRegisterReq(req_body, w, r)
	if shouldReturn {
		return
	}

	register_resp, err := ctrl.client.Do(req)
	if err != nil {
		log.Println(err.Error())
		responses.InternalError(w)
		return
	}

	defer register_resp.Body.Close()
	if register_resp.StatusCode != http.StatusOK {
		responses.ForwardResponse(w, register_resp)
	}

	err = ctrl.privileges.NewPrivilege(
		req_body.Profile.Email,
		r.Header.Get("Authorization"),
	)
	if err != nil {
		log.Println(err.Error())
		responses.InternalError(w)
		return
	}

	responses.ForwardResponse(w, register_resp)
}

func (*authCtrl) makeRegisterReq(req_body *objects.UserCreateRequest, w http.ResponseWriter, r *http.Request) (*http.Request, bool) {
	register_body, err := json.Marshal(req_body)
	if err != nil {
		log.Printf("failed to marshal register request: %s", err.Error())
		responses.ValidationErrorResponse(w, err.Error())
		return nil, true
	}
	req, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api/v1/register", utils.Config.Endpoints.IdentityProvider),
		bytes.NewBuffer(register_body),
	)
	log.Printf("Authorization %s %s", r.Header.Get("Authorization"), register_body)
	req.Header.Add("Authorization", r.Header.Get("Authorization"))
	return req, false
}

func (ctrl *authCtrl) authorize(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/authorize", utils.Config.Endpoints.IdentityProvider), r.Body)

	resp, err := ctrl.client.Do(req)
	if err != nil {
		log.Println(err.Error())
		responses.InternalError(w)
		return
	}
	if resp.StatusCode == http.StatusOK {
		data := &objects.AuthResponse{}
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, data)
		responses.JsonSuccess(w, data)
	} else {
		responses.BadRequest(w, "auth failed")
	}
}
