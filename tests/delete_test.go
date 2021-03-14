package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/klokkinn/time-service/pkg/core/models"
	"gotest.tools/v3/assert"
)

func TestDeleteEntry(t *testing.T) {
	testCases := []struct {
		name string

		withEntry models.Entry

		expectStatus int
	}{
		{
			name: "Should work",

			withEntry: models.Entry{
				Start: "1615705868",
				End:   "1615805879",
			},

			expectStatus: http.StatusNoContent,
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
			url := fmt.Sprintf("/entries/%s", id)

			deleteResult, err := env.DoRequest(url, http.MethodDelete, nil)
			assert.NilError(t, err)

			assert.Equal(t, true, bytes.Equal([]byte{}, deleteResult.Body.Bytes()))
			assert.Equal(t, tc.expectStatus, deleteResult.Code)

			getResult, err := env.DoRequest(url, http.MethodGet, nil)
			assert.NilError(t, err)

			assert.Equal(t, http.StatusNotFound, getResult.Code)
		})
	}
}
