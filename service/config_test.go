package service_test

import (
	"github.com/stone1549/yapyapyap/message/service"
	"os"
	"testing"
)

const (
	lifeCycleKey       string = "AUTH_SERVICE_ENVIRONMENT"
	repoTypeKey        string = "AUTH_SERVICE_REPO_TYPE"
	timeoutSecondsKey  string = "AUTH_SERVICE_TIMEOUT"
	portKey            string = "AUTH_SERVICE_PORT"
	pgUrlKey           string = "AUTH_SERVICE_PG_URL"
	pgInitDatasetKey   string = "AUTH_SERVICE_INIT_DATASET"
	tokenSecretKeyKey  string = "AUTH_SERVICE_TOKEN_SECRET"
	tokenPrivateKeyKey string = "AUTH_SERVICE_TOKEN_PRIV"
	tokenPublicKeyKey  string = "AUTH_SERVICE_TOKEN_PUB"
)

func clearEnv() {
	_ = os.Setenv(lifeCycleKey, "")
	_ = os.Setenv(repoTypeKey, "")
	_ = os.Setenv(timeoutSecondsKey, "")
	_ = os.Setenv(portKey, "")
	_ = os.Setenv(pgUrlKey, "")
	_ = os.Setenv(pgInitDatasetKey, "")
	_ = os.Setenv(tokenSecretKeyKey, "")
	_ = os.Setenv(tokenPrivateKeyKey, "../data/sample.key")
	_ = os.Setenv(tokenPublicKeyKey, "../data/sample.pub")
}

func setEnv(lifeCycle, repoType, timeoutSeconds, port, pgUrl, pgInitDataset, tokenSecretKey, tokenPrivateKey,
	tokenPublicKey string) {
	_ = os.Setenv(lifeCycleKey, lifeCycle)
	_ = os.Setenv(repoTypeKey, repoType)
	_ = os.Setenv(timeoutSecondsKey, timeoutSeconds)
	_ = os.Setenv(portKey, port)
	_ = os.Setenv(pgUrlKey, pgUrl)
	_ = os.Setenv(pgInitDatasetKey, pgInitDataset)
	_ = os.Setenv(tokenSecretKeyKey, tokenSecretKey)
	_ = os.Setenv(tokenPrivateKeyKey, tokenPrivateKey)
	_ = os.Setenv(tokenPublicKeyKey, tokenPublicKey)
}

// TestGetConfiguration_Defaults ensures that a default configuration is returned if no configuration is provided in
// the environment.
func TestGetConfiguration_Defaults(t *testing.T) {
	clearEnv()
	_, err := service.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_Defaults ensures that a default configuration is returned if no configuration is provided in
// the environment.
func TestGetConfiguration_ImSuccess(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "", "",
		"", "", "")
	_, err := service.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_ImSuccessSmallDataset ensures that a configuration is returned when specifying an in memory
// repo with an initial dataset.
func TestGetConfiguration_ImSuccessSmallDataset(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "",
		"../data/small_set.json", "SECRET", "", "")
	_, err := service.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_ImSuccessNoneDataset ensures that a configuration is returned when specifying an in memory
// repo without an initial dataset.
func TestGetConfiguration_ImSuccessNoneDataset(t *testing.T) {
	setEnv("DEV", "IN_MEMORY", "60", "3333", "", "",
		"SECRET", "", "")
	_, err := service.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_FailRepo ensures that an error is returned when specifying an invalid repo type.
func TestGetConfiguration_FailRepo(t *testing.T) {
	setEnv("PROD", "", "60", "3333", "", "",
		"SECRET", "", "")
	_, err := service.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_FailTimeout ensures that an error is returned when specifying an invalid timeout.
func TestGetConfiguration_FailTimeout(t *testing.T) {
	setEnv("PROD", "IN_MEMORY", "", "3333", "", "",
		"SECRET", "", "")
	_, err := service.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_FailPort ensures that an error is returned when specifying an invalid port.
func TestGetConfiguration_FailPort(t *testing.T) {
	setEnv("PROD", "IN_MEMORY", "60", "", "", "",
		"SECRET", "", "")
	_, err := service.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_PgSuccess ensures that a configuration is returned when specifying a PostgreSQL repo type.
func TestGetConfiguration_PgSuccess(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333",
		"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "",
		"SECRET", "", "")
	_, err := service.GetConfiguration()
	ok(t, err)
}

// TestGetConfiguration_PgFailPgUrl ensures that an error is returned when specifying a PostgreSQL repo type without a
// connection url.
func TestGetConfiguration_PgFailPgUrl(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333", "", "",
		"SECRET", "", "")
	_, err := service.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_ImFailNoJwtKey ensures that an error is returned when no JWT key is provided.
func TestGetConfiguration_ImFailNoJwtKey(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333", "", "",
		"", "", "")
	_, err := service.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_ImFailRsaJustPrivateKey ensures that an error is returned when no public JWT key is provided.
func TestGetConfiguration_ImFailRsaJustPrivateKey(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333", "", "",
		"", "asdfdsaf", "")
	_, err := service.GetConfiguration()
	notOk(t, err)
}

// TestGetConfiguration_ImFailRsaJustPublicKey ensures that an error is returned when no public JWT key is provided.
func TestGetConfiguration_ImFailRsaJustPublicKey(t *testing.T) {
	setEnv("PROD", "POSTGRESQL", "60", "3333", "", "",
		"", "", "asdfdsaf")
	_, err := service.GetConfiguration()
	notOk(t, err)
}
