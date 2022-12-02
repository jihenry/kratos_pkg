package cos

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"gitlab.yeahka.com/gaas/pkg/party/cos/sts"
)

type CosApi interface {
	GetTempSecret(ctx context.Context, region, bucket, path string) (*sts.CredentialResult, error)
}

type CosApiOption func(*options)

type options struct {
	secretId  string
	secretKey string
}

func WithSecretId(secretId string) CosApiOption {
	return func(opts *options) {
		if secretId != "" {
			opts.secretId = secretId
		}
	}
}

func WithSecretKey(secretKey string) CosApiOption {
	return func(opts *options) {
		if secretKey != "" {
			opts.secretKey = secretKey
		}
	}
}

type cosApiImpl struct {
	stsClient *sts.Client
}

var _ CosApi = (*cosApiImpl)(nil)

var clientMap sync.Map

func NewCosApiClient(opts ...CosApiOption) (CosApi, error) {
	options := options{}
	for _, opt := range opts {
		opt(&options)
	}
	sv, ok := clientMap.Load(fmt.Sprintf("%s:%s", options.secretId, options.secretKey))
	if ok {
		client, _ := sv.(CosApi)
		return client, nil
	}
	stsClient := sts.NewClient(
		options.secretId,
		options.secretKey,
		nil,
	)
	return &cosApiImpl{stsClient: stsClient}, nil
}

func (c *cosApiImpl) GetTempSecret(ctx context.Context, region, bucket, path string) (*sts.CredentialResult, error) {
	sfs := strings.Split(bucket, "-")
	if len(sfs) < 2 {
		return nil, fmt.Errorf("bucket:%s is invalid", bucket)
	}
	opt := &sts.CredentialOptions{
		DurationSeconds: int64(time.Hour.Seconds()),
		Region:          region,
		Policy: &sts.CredentialPolicy{
			Statement: []sts.CredentialPolicyStatement{
				{
					Action: []string{
						"name/cos:PostObject",
						"name/cos:PutObject",
					},
					Effect: "allow",
					Resource: []string{
						"qcs::cos:ap-guangzhou:uid/" + sfs[len(sfs)-1] + ":" + bucket + path,
					},
				},
			},
		},
	}
	return c.stsClient.GetCredential(opt)
}
