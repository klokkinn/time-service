package core

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"unicode/utf8"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (c Config) Validate() error {
	if c.DiscoveryURL == nil {
		return errors.New("discovery url is required")
	}

	err := is.URL.Validate(c.DiscoveryURL.String())
	if err != nil {
		return err
	}

	return validation.ValidateStruct(&c,
		validation.Field(&c.Port, validation.Required, is.UTFNumeric),
		validation.Field(&c.DSN, validation.Required),
	)
}

type KeyValueGetter func(key string) (value string)

func LoadConfig(fs *afero.Afero, kvGetter KeyValueGetter) (cfg Config) {
	getX := func(key, defaultValue string) string {
		return getValue(fs, kvGetter, key, defaultValue)
	}

	cfg.Port = getX("PORT", "3000")
	cfg.DSN = parseDSN(getX("DSN", ""))
	cfg.DiscoveryURL, _ = url.Parse(getX("DISCOVERY_URL", ""))

	switch getX("LOG_LEVEL", "info") {
	case "debug":
		cfg.LogLevel = logrus.DebugLevel
	default:
		cfg.LogLevel = logrus.InfoLevel
	}

	return cfg
}

func getValue(fs *afero.Afero, kvGetter KeyValueGetter, key, defaultValue string) (result string) {
	result = kvGetter(key)

	if result == "" {
		return defaultValue
	}

	if isPath(result) {
		var err error

		result, err = readFile(fs, result[1:])
		if err != nil {
			return defaultValue
		}
	}

	return result
}

func isPath(potentialPath string) bool {
	firstChar, _ := utf8.DecodeRuneInString(potentialPath)

	return firstChar == '@'
}

func readFile(fs *afero.Afero, path string) (result string, err error) {
	var (
		f   afero.File
		raw []byte
	)

	f, err = fs.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening file %s: %w", path, err)
	}

	defer func() {
		_ = f.Close()
	}()

	raw, err = io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("reading file %s: %w", path, err)
	}

	return string(raw), nil
}
