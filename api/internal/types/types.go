package types

type Config struct {
	Port        string
	HostName    string
	Password    string
	DBName      string
	Username    string
	SSLMode     string
	MaxOpenConn int
	Debug       bool
	MaxIdleConn int
	DB_PORT     string
}

type ClientContext string
type contextKey string

const UserContextKey contextKey = "user"
const AuthTokenKey ClientContext = "auth_token"
const DBContextKey ClientContext = "db"
const OrganizationIDKey contextKey = "organization_id"

type Payload struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
	Topic  string      `json:"topic"`
}

var JWTSecretKey = []byte("secret-key")

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
