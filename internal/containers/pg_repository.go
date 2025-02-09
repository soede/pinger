package containers

import (
	"context"
	"docker/internal/models"
)

type Repository interface {
	GetAll(ctx context.Context, page int, size int) (*models.ContainersList, error)
	GetHistory(ctx context.Context, page int, size int) (*models.ContainersList, error)
	SetAll(ctx context.Context, container []*models.Container) error
	GetByIP(ctx context.Context, q string) ([]*models.Container, error)
	GetByStatus(ctx context.Context, q string) ([]*models.Container, error)
}
