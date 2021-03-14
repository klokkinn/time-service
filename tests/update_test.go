package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/klokkinn/time-service/pkg/core/models"
	"gotest.tools/v3/assert"
)

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name string

		withEntry  *models.Entry
		withUpdate models.Entry

		expectStatus       int
		expectUpdatedEntry *models.Entry
	}{
		{
			name: "Should return not found when updating something non-existing",

			withEntry: nil,
			withUpdate: models.Entry{
				Id:  "33910bd4-ea4c-4076-a26d-4010de47334d",
				End: "1615713469",
			},

			expectStatus: http.StatusNotFound,
		},
		{
			name: "Should update with correct response including code on valid update",

			withEntry: &models.Entry{
				Start: "1615720400",
			},
			withUpdate: models.Entry{
				End: "1615720400",
			},

			expectStatus: http.StatusOK,
			expectUpdatedEntry: &models.Entry{
				Start: "1615720400",
				End:   "1615720400",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			env := createTestEnvironment(t)
			defer func() {
				_ = env.Teardown()
			}()

			update := tc.withUpdate

			if tc.withEntry == nil {
				update.Id = uuid.New().String()
			} else {
				update.Id = createEntry(t, env, *tc.withEntry)
			}

			result, err := env.DoRequest(
				fmt.Sprintf("/entries/%s", update.Id),
				http.MethodPatch,
				entryAsJSON(update),
			)
			assert.NilError(t, err)

			assert.Equal(t, tc.expectStatus, result.Code)

			if tc.expectUpdatedEntry == nil {
				return
			}

			getResult, err := env.DoRequest(fmt.Sprintf("/entries/%s", update.Id), http.MethodGet, nil)
			assert.NilError(t, err)

			tc.expectUpdatedEntry.Id = update.Id
			fetchedEntry := bytesAsEntry(getResult.Body.Bytes())

			assert.DeepEqual(t, *tc.expectUpdatedEntry, fetchedEntry)
		})
	}
}
