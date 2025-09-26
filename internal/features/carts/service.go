package carts

import (
	"context"
	"errors"

	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
	"github.com/gofiber/fiber/v2/log"
)

type ICartService interface {
	AddItem(ctx context.Context, userID string, req *CartItemRequest) error
	GetCart(ctx context.Context, userID string) ([]*CartItemsResponse, error)
	RemoveItem(ctx context.Context, userID string, productID int64) error
	ClearCart(ctx context.Context, userID string) error
}

type cartService struct {
	repo ICartRepository
}

func NewCartService(repo ICartRepository) ICartService {
	return &cartService{repo: repo}
}

func (s *cartService) AddItem(ctx context.Context, userID string, req *CartItemRequest) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	if req.Quantity <= 0 {
		return errs.ErrQuantityIsZero
	}

	err := s.repo.AddOrUpdate(ctx, userID, req.ProductID, req.Quantity)
	if err != nil {
		log.Errorf("add or update cart failed: %v", err)
		return errors.New("add items to cart failed")
	}
	return nil
}

func (s *cartService) ClearCart(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	err := s.repo.ClearCart(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrCartNotFound) {
			return err
		}
		log.Errorf("clear cart failed: %v", err)
		return errors.New("clear cart failed")
	}
	return nil
}

func (s *cartService) GetCart(ctx context.Context, userID string) ([]*CartItemsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	cart, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrCartNotFound) {
			return nil, err
		}
		log.Errorf("get cart failed: %v", err)
		return nil, errors.New("get cart failed")
	}
	return cart, nil
}

func (s *cartService) RemoveItem(ctx context.Context, userID string, productID int64) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	err := s.repo.RemoveItem(ctx, userID, productID)
	if err != nil {
		if errors.Is(err, errs.ErrCartNotFound) {
			return err
		}
		log.Errorf("remove item failed: %v", err)
		return errors.New("remove item failed")
	}
	return nil
}
