package repository

import (
	"context"
	"docker/internal/containers"
	"docker/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type containersRepo struct {
	db *sqlx.DB
}

func NewContainersRepository(db *sqlx.DB) containers.Repository {
	return &containersRepo{db: db}
}

func (r *containersRepo) GetAll(ctx context.Context, page int, size int) (*models.ContainersList, error) {
	var totalCount int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount); err != nil {
		return nil, errors.Wrap(err, "containersRepo.GetAll.QueryRowContext")
	}
	var list = make([]*models.Container, 0, totalCount)

	if err := r.db.SelectContext(ctx, &list, getAll, size, (page-1)*size); err != nil {
		return nil, errors.Wrap(err, "containersRepo.GetAll.SelectContext")
	}

	var c = &models.ContainersList{
		TotalCount: 12,
		TotalPages: 2,
		Page:       1,
		Size:       10,
		HasMore:    true, //todo
		Containers: list,
	}

	return c, nil
}
func (r *containersRepo) SetAll(ctx context.Context, containers []*models.Container) error {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "containersRepo.SetAll.BeginTxx")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.NamedExecContext(ctx, setAll, containers)
	if err != nil {
		return errors.Wrap(err, "containersRepo.SetAll.NamedExecContext")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "containersRepo.SetAll.Commit")
	}

	return nil
}

func (r *containersRepo) GetHistory(ctx context.Context, page int, size int) (*models.ContainersList, error) {
	var totalCount int
	if err := r.db.QueryRowContext(ctx, countHistoryQuery).Scan(&totalCount); err != nil {
		return nil, errors.Wrap(err, "containersRepo.GetAll.QueryRowContext")
	}
	var list = make([]*models.Container, 0, totalCount)

	if err := r.db.SelectContext(ctx, &list, getHistory, size, (page-1)*size); err != nil {
		return nil, errors.Wrap(err, "containersRepo.GetAll.SelectContext")
	}

	var c = &models.ContainersList{
		TotalCount: totalCount,
		TotalPages: 2,
		Page:       1,
		Size:       10,
		HasMore:    true, //todo
		Containers: list,
	}

	return c, nil
}

func (r *containersRepo) GetByIP(ctx context.Context, q string) ([]*models.Container, error) {
	var list []*models.Container
	err := r.db.SelectContext(ctx, &list, getByIP, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	return list, nil
}
func (r *containersRepo) GetByStatus(ctx context.Context, q string) ([]*models.Container, error) {
	var list []*models.Container
	err := r.db.SelectContext(ctx, &list, getByStatus, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	return list, nil
}
