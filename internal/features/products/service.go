package products

import (
	"context"

	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
)

type IProductService interface {
	Create(ctx context.Context, req *ProductCreate) (*Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	List(ctx context.Context, limit, offset uint) ([]*Product, error)
	Update(ctx context.Context, id int64, req *ProductUpdate) error
	Delete(ctx context.Context, id int64) error
}

type productService struct {
	repo IProductRepository
}

func NewProductService(repo IProductRepository) IProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, req *ProductCreate) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	p := &Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	}
	return s.repo.Create(ctx, p)
}

func (s *productService) GetByID(ctx context.Context, id int64) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	return s.repo.GetByID(ctx, id)
}

func (s *productService) List(ctx context.Context, limit uint, offset uint) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	if limit < 10 {
		limit = 10
	}

	return s.repo.List(ctx, limit, offset)
}

func (s *productService) Update(ctx context.Context, id int64, req *ProductUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	return s.repo.Update(ctx, id, req)
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	return s.repo.Delete(ctx, id)
}
