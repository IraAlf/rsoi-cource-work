package utils

type KafkaConfig struct {
	Endpoints []string `json:"endpoints"`
	Topics    []string `json:"topics"`
}

type EndpointConfig struct {
	Statistics       string `json:"statistics"`
	IdentityProvider string `json:"identity-provider"`
	Flights          string `json:"flights"`
	Tickets          string `json:"tickets"`
	Privileges       string `json:"privileges"`
}

type Configuration struct {
	LogFile   string         `json:"log_file"`
	Port      uint16         `json:"port"`
	RawJWKS   string         `json:"raw-jwks"`
	Kafka     KafkaConfig    `json:"kafka"`
	Endpoints EndpointConfig `json:"endpoints"`
}

var (
	Config Configuration
)

// TODO: returnable error
func InitConfig() {
	Config = Configuration{
		"logs/server.log",
		8080,
		`{"keys":[{"kty":"RSA","alg":"RS256","kid":"zavChDW9u-0gmG6lZhy4rgWspQjKQHVQ4F8hry_7ack","use":"sig","e":"AQAB","n":"kMYZfrUwocgUB-ZuzHu6qmt_Mnd4dgoOEbxLTAP-sfb2C2tBMpQ92Pa-JE7JxzDpJ485hJZh7hdObKU6cEMmLmFSCDuAfXt1dki4lhSFA8iXzRxpO6qNWkiDd48MQwLuRCC8Vog6EYGra-l3bN8j2kQR4FaK7HNDlOUWAU4qXHGhkFCEqr9rU3J74T_BcPAbGZfceyHh2a1wW84GwvGg7WYq0PmgIW5xri-VMHNNJBBIIBy9VHc84AZw0eAeNfY4G2Nf62d9Mxjs8LpSJwsYd093DealpnWapXp8ZioEiZJldEmwBtvkSI5H35upwuCABQrNFasNtno6XlmX-qw60Q"}]}`,
		KafkaConfig{Endpoints: []string{"kafka:29092"}, Topics: []string{"statistics"}},
		EndpointConfig{
			Statistics:       "http://statistics-service:8030",
			IdentityProvider: "http://identity-provider-service:8040",
			Flights:          "http://flights-service:8060",
			Tickets:          "http://tickets-service:8070",
			Privileges:       "http://privileges-service:8050",
		},
	}
}
