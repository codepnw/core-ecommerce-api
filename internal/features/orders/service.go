package orders

import (
	"context"
	"database/sql"
	"errors"

	"github.com/codepnw/core-ecommerce-system/internal/database"
	"github.com/codepnw/core-ecommerce-system/internal/features/addresses"
	"github.com/codepnw/core-ecommerce-system/internal/features/carts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
	"github.com/gofiber/fiber/v2/log"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, req *OrderRequest) error
}

type orderService struct {
	repo IOrderRepository
	cart carts.ICartService
	addr addresses.IAddressServide
	tx   *database.TxManager
}

func NewOrderService(tx *database.TxManager, repo IOrderRepository, cart carts.ICartService, addr addresses.IAddressServide) IOrderService {
	return &orderService{
		repo: repo,
		cart: cart,
		addr: addr,
		tx:   tx,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, req *OrderRequest) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Cart Total Price
	var total int64
	products, err := s.cart.GetCart(ctx, req.UserID)
	if err != nil {
		log.Errorf("get cart failed: %v", err)
		return errors.New("get cart failed")
	}

	if len(products) == 0 {
		return errors.New("cart no products")
	}

	for _, product := range products {
		subtotal := int64(product.ProductPrice) * product.ProductQuantity
		total += subtotal
	}

	// Get user address
	addr, err := s.addr.GetAddressByID(ctx, req.AddressID)
	if err != nil {
		log.Errorf("get address failed: %v", err)
		return errors.New("get address failed")
	}

	// Transaction
	err = s.tx.Transaction(ctx, func(tx *sql.Tx) error {
		// create order
		orderID, err := s.repo.InsertOrder(ctx, tx, &Order{
			UserID:     req.UserID,
			AddressID:  addr.ID,
			TotalPrice: total,
			Status:     string(StatusPending),
		})
		if err != nil {
			log.Errorf("insert order failed: %v", err)
			return errors.New("insert order failed")
		}

		// create order address
		err = s.repo.InsertOrderAddress(ctx, tx, &OrderAddress{
			OrderID:     orderID,
			AddressID:   addr.ID,
			AddressLine: addr.AddressLine,
			City:        addr.City,
			State:       addr.State,
			PostalCode:  addr.PostalCode,
			Phone:       addr.Phone,
		})
		if err != nil {
			log.Errorf("insert order_address failed: %v", err)
			return errors.New("insert order_address failed")
		}

		// create order items
		for _, product := range products {
			err = s.repo.InsertOrderItem(ctx, tx, &OrderItem{
				OrderID:   orderID,
				ProductID: product.ProductID,
				Quantity:  int(product.ProductQuantity),
				Price:     int64(product.ProductPrice),
			})
			if err != nil {
				log.Errorf("insert order_items failed: %v", err)
				return errors.New("insert order_items failed")
			}
		}

		return nil
	})
	return err
}
