package alibaba

import (
	"context"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/go-errors/errors"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/logical"
)

func (b *backend) getSTSClient(ctx context.Context, s logical.Storage, regionID string) (*sts.Client, error) {
	credential, err := b.getCredential(ctx, s)
	if err != nil {
		return nil, err
	}
	return sts.NewClientWithOptions(regionID, sdk.NewConfig(), credential)
}

func (b *backend) getRAMClient(ctx context.Context, s logical.Storage) (*ram.Client, error) {
	credential, err := b.getCredential(ctx, s)
	if err != nil {
		return nil, err
	}
	// TODO will this work? Just plugging a region because it doesn't matter?
	return ram.NewClientWithOptions("us-east-1", sdk.NewConfig(), credential)
}

func (b *backend) getECSClient(ctx context.Context, s logical.Storage, regionID string) (*ecs.Client, error) {
	credential, err := b.getCredential(ctx, s)
	if err != nil {
		return nil, err
	}
	return ecs.NewClientWithOptions(regionID, sdk.NewConfig(), credential)
}

func (b *backend) getCredential(ctx context.Context, s logical.Storage) (auth.Credential, error) {
	// TODO double-check whether you can find a way for Alibaba to use, essentially, a credential chain
	// Read the configured secret key and access key
	config, err := b.nonLockedClientConfigEntry(ctx, s)
	if err != nil {
		return nil, err
	}

	if config != nil {
		if config.AccessKey != "" && config.SecretKey != "" {
			return credentials.NewAccessKeyCredential(config.AccessKey, config.SecretKey), nil
		}
	}
	// TODO support other types of creds here
	return nil, errors.New("unable to determine credential")
}

// TODO cache clients and flush them when there's a cache clear?

// Gets an entry out of the user ID cache
func (b *backend) getCachedUserId(userId string) string {
	if userId == "" {
		return ""
	}
	if entry, ok := b.iamUserIdToArnCache.Get(userId); ok {
		b.iamUserIdToArnCache.SetDefault(userId, entry)
		return entry.(string)
	}
	return ""
}

// TODO why is it doing this?
// Sets an entry in the user ID cache
func (b *backend) setCachedUserId(userId, arn string) {
	if userId != "" {
		b.iamUserIdToArnCache.SetDefault(userId, arn)
	}
}

// TODO what is this doing?
func (b *backend) stsRoleForAccount(ctx context.Context, s logical.Storage, accountID string) (string, error) {
	// Check if an STS configuration exists for the AWS account
	sts, err := b.lockedAwsStsEntry(ctx, s, accountID)
	if err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf("error fetching STS config for account ID %q: {{err}}", accountID), err)
	}
	// An empty STS role signifies the master account
	if sts != nil {
		return sts.StsRole, nil
	}
	return "", nil
}
