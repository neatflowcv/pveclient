package fake

import (
	"context"
	"fmt"
	"strconv"

	"github.com/neatflowcv/pveclient/internal/pkg/domain"
	"github.com/neatflowcv/pveclient/internal/pkg/repository/core"
)

var _ core.CoreRepository = (*FakeCoreRepository)(nil)

type FakeCoreRepository struct {
	clusters map[string]*domain.Cluster
	nextID   int64
}

func NewFakeCoreRepository() *FakeCoreRepository {
	return &FakeCoreRepository{
		clusters: make(map[string]*domain.Cluster),
		nextID:   1,
	}
}

func (f *FakeCoreRepository) CreateCluster(ctx context.Context, spec *domain.ClusterSpec) error {
	id := strconv.FormatInt(f.nextID, 10)
	f.nextID++

	name := "cluster-" + id
	cluster := domain.NewCluster(id, spec.URL(), name)

	f.clusters[spec.URL()] = cluster

	return nil
}

func (f *FakeCoreRepository) ListClusters(ctx context.Context) ([]*domain.Cluster, error) {
	clusters := make([]*domain.Cluster, 0, len(f.clusters))
	for _, cluster := range f.clusters {
		clusters = append(clusters, cluster)
	}

	return clusters, nil
}

func (f *FakeCoreRepository) DeleteCluster(ctx context.Context, cluster *domain.Cluster) error {
	url := cluster.URL()
	if _, exists := f.clusters[url]; !exists {
		return fmt.Errorf("cluster with URL %s not found", url)
	}

	delete(f.clusters, url)

	return nil
}
