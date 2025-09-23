package products

import (
	"context"

	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
)

type IProductService interface {
	Create(ctx context.Context, req *ProductCreate) (*Product, error)
	GetByID(ctx context.Context, id int64) (*Product, error)
	List(ctx context.Context, filter *ProductFilter) ([]*Product, error)
	UpdateStock(ctx context.Context, id int64, stock int) error
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

func (s *productService) List(ctx context.Context, filter *ProductFilter) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Default Filter
	limit := 10
	offset := 0

	// Check Filter
	if filter.Limit != nil && *filter.Limit > 0 {
		limit = int(*filter.Limit)
	}

	if filter.Offset != nil && *filter.Offset > 0 {
		offset = int(*filter.Offset)
	}

	params := &ProductListParams{
		CategoryID: *filter.CategoryID,
		OrderBy:    filter.OrderBy,
		Sort:       filter.Sort,
		Limit:      limit,
		Offset:     offset,
	}

	return s.repo.List(ctx, params)
}

func (s *productService) UpdateStock(ctx context.Context, id int64, stock int) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	if stock <= 0 {
		return errs.ErrProductOutOfStock
	}

	return s.repo.UpdateStock(ctx, id, stock)
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
