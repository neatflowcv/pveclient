package flow

import (
	"context"

	"github.com/neatflowcv/pveclient/internal/pkg/domain"
	"github.com/neatflowcv/pveclient/internal/pkg/repository/core"
)

type Service struct {
	repo core.CoreRepository
}

func NewService(repo core.CoreRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCluster(ctx context.Context, spec *domain.ClusterSpec) error {
	err := s.repo.CreateCluster(ctx, spec)
	if err != nil {
		return err
	}

	return nil
}
