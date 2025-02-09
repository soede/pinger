package containers

import (
	"context"
	"docker/internal/models"
)

type UseCase interface {
	GetAll(ctx context.Context, page int, size int) (*models.ContainersList, error)
	GetHistory(ctx context.Context, page int, size int) (*models.ContainersList, error)
	GetByIP(ctx context.Context, q string) ([]*models.Container, error)
	GetByStatus(ctx context.Context, q string) ([]*models.Container, error)
	SetAll(ctx context.Context, list []*models.Container) error
}
