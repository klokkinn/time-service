/*
 * Time Service
 *
 * An API for registrating time
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"log"
	"os"

	"github.com/spf13/afero"

	"github.com/klokkinn/time-service/pkg/core"

	"github.com/klokkinn/time-service/pkg/core/router"
)

func main() {
	log.Printf("Server started")

	fs := &afero.Afero{Fs: afero.NewOsFs()}

	cfg := core.LoadConfig(fs, os.Getenv)

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	server := router.New(cfg)

	log.Fatal(server.Run(":3000"))
}
