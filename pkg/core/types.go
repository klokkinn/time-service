package core

import (
	"errors"
	"net/url"

	"github.com/klokkinn/time-service/pkg/core/models"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

/*
 * Config
 */
type Config struct {
	DiscoveryURL *url.URL
	DSN          DSN

	LogLevel logrus.Level

	Port string
}

/*
 * Storage
 */
var (
	StorageErrorNotFound = errors.New("not found")
)

type StorageFilter struct {
	Author *string
}

type StorageClient interface {
	Open() (err error)
	Close() (err error)

	Add(goal models.Entry) (err error)
	Get(id string) (result models.Entry, err error)
	GetAll(filter StorageFilter) (result []models.Entry, err error)
	Update(goal models.Entry) (updatedGoal models.Entry, err error)
	Delete(id string) (err error)
}

type DSN struct {
	Scheme       string
	Username     string
	Password     string
	URI          string
	Port         string
	DatabaseName string
}

/*
 * Router
 */
type HandlerFuncGenerator func(storage StorageClient) gin.HandlerFunc

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFuncGenerator HandlerFuncGenerator
}

// Routes is the list of the generated Route.
type Routes []Route
