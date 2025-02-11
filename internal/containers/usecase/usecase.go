package usecase

import (
	"context"
	"docker/internal/config"
	"docker/internal/containers"
	"docker/internal/models"
)

type containersUC struct {
	cfg   *config.Config
	cRepo containers.Repository
}

func NewContainersUseCase(cfg *config.Config, cRepo containers.Repository) containers.UseCase {
	return &containersUC{cfg: cfg, cRepo: cRepo}
}

func (c *containersUC) GetAll(ctx context.Context, page int, size int) (*models.ContainersList, error) {
	list, err := c.cRepo.GetAll(ctx, page, size)

	if err != nil {
		//c.logger.Error("error", err.Error())
	}
	return list, nil

}
func (c *containersUC) GetHistory(ctx context.Context, page int, size int) (*models.ContainersList, error) {
	list, err := c.cRepo.GetHistory(ctx, page, size)

	if err != nil {
		//c.logger.Error("error", err.Error())
	}
	return list, nil

}

func (c *containersUC) SetAll(ctx context.Context, list []*models.Container) error {
	err := c.cRepo.SetAll(ctx, list)

	if err != nil {
		//c.logger.Error("error", err.Error())
		return err
	}
	return nil
}

func (c *containersUC) GetByIP(ctx context.Context, q string) ([]*models.Container, error) {
	list, err := c.cRepo.GetByIP(ctx, q)
	if err != nil {
		//c.logger.Error("error", err.Error())
		return nil, err
	}
	return list, nil

}
func (c *containersUC) GetByStatus(ctx context.Context, q string) ([]*models.Container, error) {
	list, err := c.cRepo.GetByIP(ctx, q)
	if err != nil {
		//c.logger.Error("error", err.Error())
		return nil, err
	}
	return list, err
}
