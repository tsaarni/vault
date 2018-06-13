package credentials

import (
	"os"

	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/awserr"
)

// EnvProviderName provides a name of Env provider
const EnvProviderName = "EnvProvider"

var (
	// ErrAccessKeyIDNotFound is returned when the AWS Access Key ID can't be
	// found in the process's environment.
	//
	// @readonly
	ErrAccessKeyIDNotFound = awserr.New("EnvAccessKeyNotFound", "AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY not found in environment", nil)

	// ErrSecretAccessKeyNotFound is returned when the AWS Secret Access Key
	// can't be found in the process's environment.
	//
	// @readonly
	ErrSecretAccessKeyNotFound = awserr.New("EnvSecretNotFound", "AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY not found in environment", nil)
)

// A EnvProvider retrieves credentials from the environment variables of the
// running process. Environment credentials never expire.
//
// Environment variables used:
//
// * Access Key ID:     ALIBABA_ACCESS_KEY_ID
//
// * Secret Access Key: ALIBABA_ACCESS_KEY_ID
type EnvProvider struct {
	retrieved bool
}

// NewEnvCredentials returns a pointer to a new Credentials object
// wrapping the environment variable provider.
func NewEnvCredentials() *Credentials {
	return NewCredentials(&EnvProvider{})
}

// Retrieve retrieves the keys from the environment.
func (e *EnvProvider) Retrieve() (Value, error) {
	e.retrieved = false

	id := os.Getenv("ALIBABA_ACCESS_KEY_ID")
	secret := os.Getenv("ALIBABA_SECRET_ACCESS_KEY")

	if id == "" {
		return Value{ProviderName: EnvProviderName}, ErrAccessKeyIDNotFound
	}

	if secret == "" {
		return Value{ProviderName: EnvProviderName}, ErrSecretAccessKeyNotFound
	}

	e.retrieved = true
	return Value{
		AccessKeyID:     id,
		SecretAccessKey: secret,
		SessionToken:    os.Getenv("ALIBABA_SESSION_TOKEN"),
		ProviderName:    EnvProviderName,
	}, nil
}

// IsExpired returns if the credentials have been retrieved.
func (e *EnvProvider) IsExpired() bool {
	return !e.retrieved
}
