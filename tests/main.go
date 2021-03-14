package tests

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/klokkinn/time-service/pkg/core"
	"github.com/klokkinn/time-service/pkg/core/models"
	"github.com/klokkinn/time-service/pkg/core/router"
	"github.com/oslokommune/go-gin-tools/pkg/v1/servicetesting"
	authtesting "github.com/oslokommune/go-oidc-middleware/pkg/v1/testing"
	"gotest.tools/assert"
)

func createTestEnvironment(t *testing.T) *servicetesting.Environment {
	authTestOptions := authtesting.NewTestTokenOptions()

	discoveryServer := authtesting.CreateTestDiscoveryServer(authTestOptions)
	bearerToken := authtesting.CreateTestToken(authTestOptions)

	discoveryURL, _ := url.Parse(discoveryServer.URL)
	dbPassword := "postgres"

	env, err := servicetesting.NewGinTestEnvironment(servicetesting.CreatePostgresDatabaseBackendOptions(dbPassword), bearerToken)
	assert.NilError(t, err)

	parts := strings.Split(env.GetDatabaseBackendURI(), ":")
	dbURI := parts[0]
	dbPort := parts[1]

	cfg := core.Config{
		DiscoveryURL: discoveryURL,
		DSN: core.DSN{
			Scheme:       "postgres",
			Username:     "postgres",
			Password:     dbPassword,
			URI:          dbURI,
			Port:         dbPort,
			DatabaseName: "postgres",
		},
		Port: "3000",
	}

	err = cfg.Validate()
	if err != nil {
		_ = env.Teardown()

		assert.NilError(t, err)
	}

	env.TestServer = router.New(cfg)

	return env
}

func entryAsJSON(entry models.Entry) (raw []byte) {
	raw, _ = json.Marshal(entry)

	return raw
}

func bytesAsEntry(raw []byte) (entry models.Entry) {
	_ = json.Unmarshal(raw, &entry)

	return entry
}

func createEntry(t *testing.T, env *servicetesting.Environment, entry models.Entry) (id string) {
	result, err := env.DoRequest("/entries", http.MethodPost, entryAsJSON(entry))
	assert.NilError(t, err)

	var createdEntry models.Entry

	err = json.Unmarshal(result.Body.Bytes(), &createdEntry)
	assert.NilError(t, err)

	return createdEntry.Id
}

//func createEntries(t *testing.T, env *servicetesting.Environment, entries []models.Entry) (ids []string) {
//	ids = make([]string, len(entries))
//
//	for index, entry := range entries {
//		id := createEntry(t, env, entry)
//
//		ids[index] = id
//	}
//
//	return ids
//}
