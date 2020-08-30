package env

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aaronsky/applereleaser/pkg/context"
)

// ErrMissingEnvVar indicates an error when a required variable is missing in the environment.
var ErrMissingEnvVar = errors.New("missing required environment variable")

// Pipe is a global hook pipe.
type Pipe struct{}

// String is the name of this pipe.
func (Pipe) String() string {
	return "loading environment variables"
}

// Run executes the hooks.
func (p Pipe) Run(ctx *context.Context) error {
	keyID, err := loadEnv("ASC_KEY_ID", true)
	if err != nil {
		return err
	}
	issuerID, err := loadEnv("ASC_ISSUER_ID", true)
	if err != nil {
		return err
	}
	privateKey, err := loadEnvFromPath("ASC_PRIVATE_KEY_PATH", true)
	if err != nil {
		return err
	}
	ctx.Credentials.KeyID = keyID
	ctx.Credentials.IssuerID = issuerID
	ctx.Credentials.PrivateKey = privateKey
	return nil
}

func loadEnv(env string, required bool) (string, error) {
	val := os.Getenv(env)
	if val == "" && required {
		return "", fmt.Errorf("key %s not found: %w", env, ErrMissingEnvVar)
	}
	return val, nil
}

func loadEnvFromPath(env string, required bool) (string, error) {
	val, err := loadEnv(env, required)
	if err != nil {
		return "", err
	}
	f, err := os.Open(val)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(f)
	return string(bytes), err
}