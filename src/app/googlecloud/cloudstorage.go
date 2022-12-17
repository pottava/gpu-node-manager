package googlecloud

import (
	"context"
	"strings"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/storage"
	"github.com/pottava/gpu-node-manager/src/app/util"
)

func MakeBucket(ctx context.Context, name string) error {
	client, err := storage.NewClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	err = client.Bucket(name).Create(ctx, util.ProjectID(), &storage.BucketAttrs{
		Location: util.Location,
		Autoclass: &storage.Autoclass{
			Enabled: true,
		},
		PublicAccessPrevention: storage.PublicAccessPreventionEnforced,
		UniformBucketLevelAccess: storage.UniformBucketLevelAccess{
			Enabled: true,
		},
	})
	if err != nil && strings.Contains(err.Error(), "you already own it") {
		return nil
	}
	return err
}

func AddRoleToBucket(ctx context.Context, name, email, role string) error {
	client, err := storage.NewClient(ctx, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	policy := &iam.Policy{}
	policy.Add(email, iam.RoleName("roles/"+role))
	return client.Bucket(name).IAM().SetPolicy(ctx, policy)
}
