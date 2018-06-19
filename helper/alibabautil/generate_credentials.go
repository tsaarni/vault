package alibabautil

import (
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
)

type CredentialsConfig struct {
	// The access key if static credentials are being used
	AccessKey string

	// The secret key if static credentials are being used
	SecretKey string

	// The session token if it is being used
	SessionToken string

	// If specified, the region will be provided to the config of the
	// EC2RoleProvider's client. This may be useful if you want to e.g. reuse
	// the client elsewhere.
	Region string

	// The filename for the shared credentials provider, if being used
	Filename string

	// The profile for the shared credentials provider, if being used
	Profile string

	// The http.Client to use, or nil for the client to use its default
	HTTPClient *http.Client
}

// TODO - support more types of credentials,
// and strip params above for ones you won't support
// and consider pulling them from the env to be more aws-y or a file
// although the sdk doesn't support it.
// Or don't because security risk?
func (c *CredentialsConfig) GenerateCredentialChain() *credentials.AccessKeyCredential {
	return &credentials.AccessKeyCredential{
		AccessKeyId:     c.AccessKey,
		AccessKeySecret: c.SecretKey,
	}
}
