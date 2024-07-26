package utils

type DatabaseConfiguration struct {
	Type string `json:"type"`
	Name string `json:"name"`

	User     string `json:"user"`
	Password string `json:"password"`

	Host string `json:"host"`
	Port string `json:"port"`
}

type Configuration struct {
	DB      DatabaseConfiguration `json:"db"`
	LogFile string                `json:"logFile"`
	Port    uint16                `json:"port"`
	RawJWKS string                `json:"rawJWKS"`
}

var Config Configuration

func InitConfiguration() {
	Config = Configuration{
		DatabaseConfiguration{
			"postgres",
			"identity_provider",
			"program",
			"test",
			"postgres-service",
			"5432",
		},

		"logs/servel.log",
		8040,
		`{"keys":[{"kty":"RSA","alg":"RS256","kid":"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsBCnQuiAgfjKlTnTdfH+ZTFaic0aqI9/7zd0S2Fyi+JqeVNkmLYDBEvl6t+Ju8WLY3P55dNq6lvPpskq69GrarHmjxiti2LqZoay1rUqhccFK/tMuHENQV4Wqiy7FMTIGicZ+qq1mAiUZo7jU0YBnZTOdIwF7p1CDm1QJcxsvhaU41Q9zBq4N3BMfP2wwBPyat6cPqbh3ok98KmXIR5CGw5EdlsImW4EXU2v7f8ZGpTVb4QqQip0F+jKqadG6xIvKyaSNk7tThMoN7DAMMXVeXnoydVIdesZE0E3a0LzVU4QFImhKsULC73wcP8+jYxlTCj8XTE2Ruvk8F+KpYct+QIDAQAB","use":"sig","e":"AQAB","n":"kMYZfrUwocgUB-ZuzHu6qmt_Mnd4dgoOEbxLTAP-sfb2C2tBMpQ92Pa-JE7JxzDpJ485hJZh7hdObKU6cEMmLmFSCDuAfXt1dki4lhSFA8iXzRxpO6qNWkiDd48MQwLuRCC8Vog6EYGra-l3bN8j2kQR4FaK7HNDlOUWAU4qXHGhkFCEqr9rU3J74T_BcPAbGZfceyHh2a1wW84GwvGg7WYq0PmgIW5xri-VMHNNJBBIIBy9VHc84AZw0eAeNfY4G2Nf62d9Mxjs8LpSJwsYd093DealpnWapXp8ZioEiZJldEmwBtvkSI5H35upwuCABQrNFasNtno6XlmX-qw60Q"}]}`,
	}
}
