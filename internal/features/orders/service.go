package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/codepnw/core-ecommerce-system/internal/database"
	"github.com/codepnw/core-ecommerce-system/internal/features/addresses"
	"github.com/codepnw/core-ecommerce-system/internal/features/carts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
	"github.com/gofiber/fiber/v2/log"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, req *OrderRequest) error
	ListOrders(ctx context.Context, filter *OrderFilter) ([]*OrdersResponse, error)
	UpdateOrderStatus(ctx context.Context, id int64, status OrderStatus) error
}

type OrderServiceConfig struct {
	OrderRepo IOrderRepository          `validate:"required"`
	CartSrv   carts.ICartService        `validate:"required"`
	AddrSrv   addresses.IAddressServide `validate:"required"`
	Tx        *database.TxManager       `validate:"required"`
}

func NewOrderService(cfg *OrderServiceConfig) (IOrderService, error) {
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("OrderServiceConfig required all fields: %w", err)
	}
	return cfg, nil
}

func (s *OrderServiceConfig) CreateOrder(ctx context.Context, req *OrderRequest) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// Cart Total Price
	var total int64
	products, err := s.CartSrv.GetCart(ctx, req.UserID)
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
	addr, err := s.AddrSrv.GetAddressByID(ctx, req.AddressID)
	if err != nil {
		log.Errorf("get address failed: %v", err)
		return errors.New("get address failed")
	}

	// Transaction
	err = s.Tx.Transaction(ctx, func(tx *sql.Tx) error {
		// create order
		orderID, err := s.OrderRepo.InsertOrder(ctx, tx, &Order{
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
		err = s.OrderRepo.InsertOrderAddress(ctx, tx, &OrderAddress{
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
			err = s.OrderRepo.InsertOrderItem(ctx, tx, &OrderItem{
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

func (s *OrderServiceConfig) ListOrders(ctx context.Context, filter *OrderFilter) ([]*OrdersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	res, err := s.OrderRepo.ListOrders(ctx, filter)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *OrderServiceConfig) UpdateOrderStatus(ctx context.Context, id int64, status OrderStatus) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	return s.OrderRepo.UpdateStatus(ctx, id, string(status))
}
