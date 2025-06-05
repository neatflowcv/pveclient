package core

import (
	"context"

	"github.com/neatflowcv/pveclient/internal/pkg/domain"
)

type CoreRepository interface {
	CreateCluster(ctx context.Context, spec *domain.ClusterSpec) error
	ListClusters(ctx context.Context) ([]*domain.Cluster, error)
	DeleteCluster(ctx context.Context, cluster *domain.Cluster) error
}
