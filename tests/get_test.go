package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/klokkinn/time-service/pkg/core/models"
	"gotest.tools/v3/assert"
)

func TestGetEntry(t *testing.T) {
	testCases := []struct {
		name string

		withEntry models.Entry

		expectStatus int
	}{
		{
			name: "Should return an entry after post with the correct status",

			withEntry: models.Entry{Start: "1615705868", End: "1615705879"},

			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			env := createTestEnvironment(t)
			defer func() {
				_ = env.Teardown()
			}()

			id := createEntry(t, env, tc.withEntry)

			result, err := env.DoRequest(fmt.Sprintf("/entries/%s", id), http.MethodGet, nil)
			assert.NilError(t, err)

			fetchedEvent := bytesAsEntry(result.Body.Bytes())
			fetchedEvent.Id = ""

			assert.DeepEqual(t, tc.withEntry, fetchedEvent)
		})
	}
}

func TestMissingEntry(t *testing.T) {
	testCases := []struct {
		name string

		withId string

		expectStatus int
	}{
		{
			name: "Should work",

			withId: "ea151b9f-1d4f-4800-8dbd-c2a47bec1f45",

			expectStatus: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			env := createTestEnvironment(t)
			defer func() {
				_ = env.Teardown()
			}()

			result, err := env.DoRequest(fmt.Sprintf("/entries/%s", tc.withId), http.MethodGet, nil)
			assert.NilError(t, err)

			assert.Equal(t, tc.expectStatus, result.Code)
		})
	}
}
