/*
 * Time Service
 *
 * An API for registrating time
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package handlers

import (
	"errors"
	"net/http"

	"github.com/klokkinn/time-service/pkg/core/models"

	"github.com/klokkinn/time-service/pkg/core"

	"github.com/gin-gonic/gin"
)

// DeleteEntry - Delete an existing time Entry
func DeleteEntry(storage core.StorageClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err error
			id  string
		)

		id = c.Param("id")
		if id == "" {
			c.Status(http.StatusBadRequest)

			return
		}

		err = storage.Delete(id)
		if err != nil {
			var storageError *core.StorageError

			if errors.As(err, &storageError) {
				c.Status(storageError.StatusCode())
			} else {
				c.Status(http.StatusInternalServerError)
			}

			return
		}

		c.Status(http.StatusNoContent)
	}
}

// GetEntry - Get detailed information about a single time Entry
func GetEntry(storage core.StorageClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err   error
			id    string
			entry models.Entry
		)

		id = c.Param("id")
		if id == "" {
			c.Status(http.StatusBadRequest)

			return
		}

		entry, err = storage.Get(id)
		if err != nil {
			var storageError *core.StorageError

			if errors.As(err, &storageError) {
				c.Status(storageError.StatusCode())
			} else {
				c.Status(http.StatusInternalServerError)
			}

			return
		}

		c.JSON(http.StatusOK, entry)
	}
}

// UpdateEntry - Edit an existing time Entry
func UpdateEntry(storage core.StorageClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err   error
			id    string
			entry models.Entry
		)

		id = c.Param("id")
		if id == "" {
			c.Status(http.StatusBadRequest)

			return
		}

		err = c.Bind(&entry)
		if err != nil {
			c.Status(http.StatusBadRequest)

			return
		}

		entry, err = storage.Update(entry)
		if err != nil {
			var storageError *core.StorageError

			if errors.As(err, &storageError) {
				c.Status(storageError.StatusCode())
			} else {
				c.Status(http.StatusInternalServerError)
			}

			return
		}

		c.JSON(http.StatusOK, entry)
	}
}
