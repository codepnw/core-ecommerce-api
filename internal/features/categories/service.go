package categories

import (
	"context"

	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
)

type CategoryService interface {
	Create(ctx context.Context, req *CategoryCreate) error
	List(ctx context.Context) ([]*Category, error)
	Update(ctx context.Context, id int64, req *CategoryUpdate) error
	Delete(ctx context.Context, id int64) error
}

type categoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, req *CategoryCreate) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	input := &Category{
		Name:        req.Name,
		Description: req.Description,
	}
	return s.repo.Create(ctx, input)
}

func (s *categoryService) List(ctx context.Context) ([]*Category, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	return s.repo.List(ctx)
}

func (s *categoryService) Update(ctx context.Context, id int64, req *CategoryUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	if req.Name == nil && req.Description == nil {
		return errs.ErrNoFieldUpdate
	}

	input := &Category{
		ID:          id,
		Name:        *req.Name,
		Description: *req.Description,
	}
	return s.repo.Update(ctx, input)
}

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	return s.repo.Delete(ctx, id)
}
