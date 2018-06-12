package sts

import (
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws"
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/client"
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/client/metadata"
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/request"
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/aws/signer/v4"
	"github.com/hashicorp/vault/builtin/credential/alibaba/alibaba-sdk-go/private/protocol/query"
)

// STS provides the API operation methods for making requests to
// AWS Security Token Service. See this package's package overview docs
// for details on the service.
//
// STS methods are safe to use concurrently. It is not safe to
// modify mutate any of the struct's properties though.
type STS struct {
	*client.Client
}

// Used for custom client initialization logic
var initClient func(*client.Client)

// Used for custom request initialization logic
var initRequest func(*request.Request)

// Service information constants
const (
	ServiceName = "sts"       // Service endpoint prefix API calls made to.
	EndpointsID = ServiceName // Service ID for Regions and Endpoints metadata.
)

// New creates a new instance of the STS client with a session.
// If additional configuration is needed for the client instance use the optional
// aws.Config parameter to add your extra config.
//
// Example:
//     // Create a STS client from just a session.
//     svc := sts.New(mySession)
//
//     // Create a STS client with additional configuration
//     svc := sts.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
func New(p client.ConfigProvider, cfgs ...*aws.Config) (*STS, error) {
	c := p.ClientConfig(EndpointsID, cfgs...)
	return newClient(*c.Config, c.Handlers, c.Endpoint, c.SigningRegion, c.SigningName)
}

// newClient creates, initializes and returns a new service client instance.
func newClient(cfg aws.Config, handlers request.Handlers, endpoint, signingRegion, signingName string) (*STS, error) {
	svc := &STS{
		Client: client.New(
			cfg,
			metadata.ClientInfo{
				ServiceName:   ServiceName,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				Endpoint:      endpoint,
				APIVersion:    "2011-06-15",
			},
			handlers,
		),
	}

	creds, err := cfg.Credentials.Get()
	if err != nil {
		return nil, err
	}

	// Handlers
	svc.Handlers.Build.PushBackNamed(query.NewBuildHandler(creds.AccessKeyID))
	svc.Handlers.Sign.PushBackNamed(v4.NewSignRequestHandler(creds.SecretAccessKey))
	svc.Handlers.Unmarshal.PushBackNamed(query.UnmarshalHandler)
	svc.Handlers.UnmarshalMeta.PushBackNamed(query.UnmarshalMetaHandler)
	svc.Handlers.UnmarshalError.PushBackNamed(query.UnmarshalErrorHandler)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc.Client)
	}

	return svc, nil
}

// newRequest creates a new request for a STS operation and runs any
// custom request initialization.
func (c *STS) newRequest(op *request.Operation, params, data interface{}) *request.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(req)
	}

	return req
}
