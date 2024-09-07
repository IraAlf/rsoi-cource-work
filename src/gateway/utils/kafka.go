package utils

import (
	"encoding/json"
	"gateway/objects"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/golang-jwt/jwt"
	"github.com/urfave/negroni"
)

func TokenIsMissing(w http.ResponseWriter) {
	msg := "Missing auth token"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(msg)
}

func JwtAccessDenied(w http.ResponseWriter) {
	msg := "jwt-token is not valid"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(msg)
}

func TokenExpired(w http.ResponseWriter) {
	msg := "jwt-token expired"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(msg)
}

func ForwardResponse(w http.ResponseWriter, resp *http.Response) {
	w.WriteHeader(resp.StatusCode)
	body := []byte{}
	if resp.ContentLength != 0 {
		body, _ = ioutil.ReadAll(resp.Body)
	}
	json.NewEncoder(w).Encode(body)
}

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

func RetrieveToken(w http.ResponseWriter, r *http.Request) *Claims {
	reqToken := r.Header.Get("Authorization")
	log.Printf("Privileges: token: %s ", reqToken)
	if len(reqToken) == 0 {
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
		return nil
	}
	if time.Now().Unix()-tk.ExpiresAt > 0 {
		log.Printf("ExpiresAt")
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

func sendRequestStatToKafka(stat *objects.RequestStat, topic string, producer sarama.SyncProducer) {
	statBytes, err := json.Marshal(stat)
	if err != nil {
		log.Printf("Error encoding request stats")
		return
	}

	// Создаем сообщение Kafka
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(statBytes),
	}

	// Отправляем сообщение в Kafka
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Printf("Error sending request stat to Kafka: %v", err)
		return
	}

	log.Printf("Request stat sent to Kafka: %s", string(statBytes))
}

// Обертка для обработчиков HTTP, чтобы сохранять статистику запросов
func RequestStatMiddleware(next http.Handler, topic string, producer sarama.SyncProducer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := negroni.NewResponseWriter(w)
		token := RetrieveToken(w, r)

		stat := &objects.RequestStat{}
		stat.Path = r.URL.Path
		stat.Method = r.Method
		if token != nil {
			stat.UserName = token.Subject
		} else {
			stat.UserName = ""
		}
		log.Printf("%s", r.Header.Get("X-User-Name"))
		log.Printf("User name: %s", stat.UserName)
		stat.StartedAt = time.Now()

		// Вывод всех заголовков
		log.Println("Headers:")
		for name, values := range r.Header {
			for _, value := range values {
				log.Printf("%s: %s", name, value)
			}
		}

		next.ServeHTTP(lrw, r)

		stat.FinishedAt = time.Now()
		stat.Duration = stat.FinishedAt.Sub(stat.StartedAt)
		stat.ResponceCode = lrw.Status()

		go sendRequestStatToKafka(stat, topic, producer)
	})
}
