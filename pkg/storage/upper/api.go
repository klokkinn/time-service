package upper

import (
	"fmt"

	"github.com/klokkinn/time-service/pkg/core"

	"github.com/upper/db/v4/adapter/postgresql"
)

func NewClient(dsn core.DSN) core.StorageClient {
	return &client{
		connectionURL: &postgresql.ConnectionURL{
			User:     dsn.Username,
			Password: dsn.Password,
			Host:     fmt.Sprintf("%s:%s", dsn.URI, dsn.Port),
			Database: dsn.DatabaseName,
		},
	}
}
