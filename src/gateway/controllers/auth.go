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

var jwtKey = []byte("your-256-bit-secret")

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
	log.Printf("Flights: token: %s ", reqToken)
	if len(reqToken) == 0 {
		responses.TokenIsMissing(w)
		return nil, nil
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
		return nil, nil
	}
	if time.Now().Unix()-tk.ExpiresAt > 0 {
		log.Printf("ExpiresAt")
		responses.TokenExpired(w)
		return nil, nil
	}
	log.Printf("Flights: token: %s ", tk)
	return tk, nil
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
		fmt.Println("USERTYPE\n\n\n %s", data.Role)
		responses.JsonSuccess(w, data)
	} else {
		responses.BadRequest(w, "auth failed")
	}
}
