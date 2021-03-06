package hub

import (
	"fmt"
	"os"
	"strings"
)

// Options stores the hub's options
type Options struct {
	Debug                 bool
	DBPath                string
	PublisherJWTKey       []byte
	SubscriberJWTKey      []byte
	AllowAnonymous        bool
	CorsAllowedOrigins    []string
	PublishAllowedOrigins []string
	Addr                  string
	AcmeHosts             []string
	AcmeCertDir           string
	CertFile              string
	KeyFile               string
	Demo                  bool
}

func getJWTKey(role string) string {
	key := os.Getenv(fmt.Sprintf("%s_JWT_KEY", role))
	if key == "" {
		return os.Getenv("JWT_KEY")
	}

	return key
}

// NewOptionsFromEnv creates a new option instance from environment
// It returns an error if mandatory env env vars are missing
func NewOptionsFromEnv() (*Options, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "updates.db"
	}

	options := &Options{
		os.Getenv("DEBUG") == "1",
		dbPath,
		[]byte(getJWTKey("PUBLISHER")),
		[]byte(getJWTKey("SUBSCRIBER")),
		os.Getenv("ALLOW_ANONYMOUS") == "1",
		splitVar(os.Getenv("CORS_ALLOWED_ORIGINS")),
		splitVar(os.Getenv("PUBLISH_ALLOWED_ORIGINS")),
		os.Getenv("ADDR"),
		splitVar(os.Getenv("ACME_HOSTS")),
		os.Getenv("ACME_CERT_DIR"),
		os.Getenv("CERT_FILE"),
		os.Getenv("KEY_FILE"),
		os.Getenv("DEMO") == "1" || os.Getenv("DEBUG") == "1",
	}

	missingEnv := make([]string, 0, 4)
	if len(options.PublisherJWTKey) == 0 {
		missingEnv = append(missingEnv, "PUBLISHER_JWT_KEY")
	}
	if len(options.SubscriberJWTKey) == 0 {
		missingEnv = append(missingEnv, "SUBSCRIBER_JWT_KEY")
	}
	if len(options.CertFile) != 0 && len(options.KeyFile) == 0 {
		missingEnv = append(missingEnv, "KEY_FILE")
	}
	if len(options.KeyFile) != 0 && len(options.CertFile) == 0 {
		missingEnv = append(missingEnv, "CERT_FILE")
	}

	if len(missingEnv) > 0 {
		return nil, fmt.Errorf("The following environment variable must be defined: %s", missingEnv)
	}

	return options, nil
}

func splitVar(v string) []string {
	if v == "" {
		return []string{}
	}

	return strings.Split(v, ",")
}
