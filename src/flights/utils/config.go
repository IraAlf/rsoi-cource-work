package utils

type DBConfiguration struct {
	Type string `json:"type"`
	Name string `json:"name"`

	User     string `json:"user"`
	Password string `json:"password"`

	Host string `json:"host"`
	Port string `json:"port"`
}

type Configuration struct {
	DB      DBConfiguration `json:"db"`
	LogFile string          `json:"log_file"`
	Port    uint16          `json:"port"`
	RawJWKS string          `json:"raw-jwks"`
}

var (
	Config Configuration
)

// TODO: returnable error
func InitConfig() {
	Config = Configuration{
		DBConfiguration{
			"postgres",
			"flights",
			"program",
			"test",
			"postgres-service",
			"5432",
		},
		"logs/server.log",
		8060,
		`{"keys":[{"kty":"RSA","alg":"RS256","kid":"your-key-id"}]}`,
	}
}
