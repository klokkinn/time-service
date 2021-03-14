package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/klokkinn/time-service/pkg/core/models"
	"gotest.tools/v3/assert"
)

func TestAddEntry(t *testing.T) {
	testCases := []struct {
		name string

		withEntry models.Entry

		expectStatus int
	}{
		{
			name: "Should add entry and return correct status upon success",

			withEntry: models.Entry{
				Start: "1615705868",
				End:   "1615805879",
			},

			expectStatus: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			env := createTestEnvironment(t)
			defer func() {
				_ = env.Teardown()
			}()

			postResult, err := env.DoRequest("/entries", http.MethodPost, entryAsJSON(tc.withEntry))

			assert.NilError(t, err)
			assert.Equal(t, tc.expectStatus, postResult.Code)

			id := bytesAsEntry(postResult.Body.Bytes()).Id

			getResult, err := env.DoRequest(fmt.Sprintf("/entries/%s", id), http.MethodGet, nil)
			assert.NilError(t, err)

			fetchedEntry := bytesAsEntry(getResult.Body.Bytes())
			fetchedEntry.Id = ""

			assert.DeepEqual(t, tc.withEntry, fetchedEntry)
		})
	}
}
