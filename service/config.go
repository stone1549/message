package service

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	lifeCycleKey      string = "MESSAGE_SERVICE_ENVIRONMENT"
	repoTypeKey       string = "MESSAGE_SERVICE_REPO_TYPE"
	timeoutSecondsKey string = "MESSAGE_SERVICE_TIMEOUT"
	portKey           string = "MESSAGE_SERVICE_PORT"
	pgUrlKey          string = "MESSAGE_SERVICE_PG_URL"
	initDatasetKey    string = "MESSAGE_SERVICE_INIT_DATASET"
	tokenSecretKeyKey string = "MESSAGE_SERVICE_TOKEN_SECRET"
	tokenPrivateKey   string = "MESSAGE_SERVICE_TOKEN_PRIV"
	tokenPublicKey    string = "MESSAGE_SERVICE_TOKEN_PUB"
)

// LifeCycle represents a particular application life cycle.
type LifeCycle int

const (
	// DevLifeCycle represents the development environment.
	DevLifeCycle LifeCycle = 0
	// PreProdLifeCycle represents the pre production environment.
	PreProdLifeCycle LifeCycle = iota
	// ProdLifeCycle represents the production environment.
	ProdLifeCycle LifeCycle = iota
)

func (lc LifeCycle) String() string {
	switch lc {
	case DevLifeCycle:
		return "DEV"
	case PreProdLifeCycle:
		return "PRE_PROD"
	case ProdLifeCycle:
		return "PROD"
	default:
		return ""
	}
}

// MessageRepositoryType represents a type of UserRepository
type MessageRepositoryType int

const (
	// InMemoryRepo represents a UserRepository that is stored entirely in memory.
	InMemoryRepo MessageRepositoryType = 0
	// PostgreSqlRepo represents a UserRepository that utilizes a PostgreSQL database.
	PostgreSqlRepo MessageRepositoryType = iota
)

func (prt MessageRepositoryType) String() string {
	switch prt {
	case PostgreSqlRepo:
		return "POSTGRESQL"
	case InMemoryRepo:
		return "IN_MEMORY"
	default:
		return ""
	}
}

// Configuration provides methods for retrieving aspects of the application's configuration.
type Configuration interface {
	// GetLifeCycle retrieves the configured life cycle.
	GetLifeCycle() LifeCycle
	// GetRepoType retrieves the configured repo type.
	GetRepoType() MessageRepositoryType
	// GetTimeout retrieves the configured request timeout.
	GetTimeout() time.Duration
	// GetPort retrieves the configured port.
	GetPort() int

	// GetInitDataSet retrieves the path to an initial dataset to load on app launch, mostly for testing and dev use.
	GetInitDataSet() string

	// GetPgUrl retrieves the configured url string for connecting to PostgreSQL.
	GetPgUrl() string

	// GetTokenSecretKey a shared secret key for signing tokens
	GetTokenSecretKey() string

	// GetTokenPrivateKey retrieves the the private key used to sign tokens.
	GetTokenPrivateKey() *rsa.PrivateKey

	// GetTokenPublicKey retrieves public key used to validate JWT tokens.
	GetTokenPublicKey() *rsa.PublicKey
}

type configuration struct {
	lifeCycle   LifeCycle
	repoType    MessageRepositoryType
	timeout     time.Duration
	port        int
	pgUrl       string
	initDataset string
	secretKey   string
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
}

func (conf *configuration) GetLifeCycle() LifeCycle {
	return conf.lifeCycle
}

func (conf *configuration) GetRepoType() MessageRepositoryType {
	return conf.repoType
}

func (conf *configuration) GetTimeout() time.Duration {
	return conf.timeout
}

func (conf *configuration) GetPort() int {
	return conf.port
}

func (conf *configuration) GetPgUrl() string {
	return conf.pgUrl
}

func (conf *configuration) GetInitDataSet() string {
	return conf.initDataset
}

// GetTokenSecretKey retrieves the shared secret key for reading/signing JWT tokens
func (conf *configuration) GetTokenSecretKey() string {
	return conf.secretKey
}

// GetTokenPrivateKey retrieves the private RSA key used to sign JWT tokens.
func (conf *configuration) GetTokenPrivateKey() *rsa.PrivateKey {
	return conf.privateKey
}

// GetTokenPublicKey retrieves public RSA key used to validate JWT tokens.
func (conf *configuration) GetTokenPublicKey() *rsa.PublicKey {
	return conf.publicKey
}

// GetConfiguration constructs a Configuration based on environment variables.
func GetConfiguration() (Configuration, error) {
	var err error
	config := configuration{}

	lcStr := os.Getenv(lifeCycleKey)

	switch lcStr {
	case DevLifeCycle.String():
		config.lifeCycle = DevLifeCycle
	case PreProdLifeCycle.String():
		config.lifeCycle = PreProdLifeCycle
	case ProdLifeCycle.String():
		config.lifeCycle = ProdLifeCycle
	default:
		config.lifeCycle = DevLifeCycle
	}

	if err != nil {
		return nil, err
	}

	repoTypeStr := os.Getenv(repoTypeKey)

	switch repoTypeStr {
	case InMemoryRepo.String():
		config.repoType = InMemoryRepo
	case PostgreSqlRepo.String():
		config.repoType = PostgreSqlRepo
	default:
		if config.lifeCycle == DevLifeCycle {
			config.repoType = InMemoryRepo
		} else {
			err = errors.New(fmt.Sprintf("No repo type configured, set %s environment variable", repoTypeKey))
		}
	}

	if err != nil {
		return nil, err
	}

	timeoutStr := os.Getenv(timeoutSecondsKey)

	if timeoutStr == "" && config.lifeCycle == DevLifeCycle {
		timeoutStr = "60"
	}

	timeoutInt, err := strconv.Atoi(timeoutStr)

	if err != nil {
		err = errors.New(fmt.Sprintf("No timeout configured, set %s environment variable", timeoutSecondsKey))
		return nil, err
	}

	config.timeout = time.Duration(timeoutInt) * time.Second

	portStr := os.Getenv(portKey)

	if portStr == "" && config.lifeCycle == DevLifeCycle {
		portStr = "3333"
	}
	port, err := strconv.Atoi(portStr)

	if err != nil {
		err = errors.New(fmt.Sprintf("No port configured, set %s environment variable", portKey))
		return nil, err
	}

	config.port = port

	if config.repoType == PostgreSqlRepo {
		err = setPostgresqlConfig(&config)
	}

	if err != nil {
		return nil, err
	}

	config.initDataset = os.Getenv(initDatasetKey)

	if err != nil {
		return nil, err
	}

	secretKey := os.Getenv(tokenSecretKeyKey)
	privateKeyPath := os.Getenv(tokenPrivateKey)
	publicKeyPath := os.Getenv(tokenPublicKey)

	if secretKey == "" && (privateKeyPath == "" || publicKeyPath == "") {
		if config.lifeCycle == DevLifeCycle {
			config.secretKey = "secret"
		} else {
			return nil, errors.New(fmt.Sprintf("must set either %s environment variable or both %s AND %s "+
				"environment variables", tokenSecretKeyKey, tokenPublicKey, tokenPrivateKey))
		}
	} else if secretKey != "" {
		config.secretKey = secretKey
	} else {
		signBytes, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, err
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
		if err != nil {
			return nil, err
		}

		verifyBytes, err := os.ReadFile(publicKeyPath)
		if err != nil {
			return nil, err
		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

		config.privateKey = privateKey
		config.publicKey = publicKey
	}

	return &config, nil
}

func setPostgresqlConfig(config *configuration) error {
	var err error

	config.pgUrl = os.Getenv(pgUrlKey)

	if strings.TrimSpace(config.pgUrl) == "" {
		err = errors.New(fmt.Sprintf("No PostgreSqlRepo url configured, set %s environment variable", pgUrlKey))
	}

	return err
}
