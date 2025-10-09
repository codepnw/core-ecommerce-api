package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/codepnw/core-ecommerce-system/internal/database"
	"github.com/codepnw/core-ecommerce-system/internal/features/addresses"
	"github.com/codepnw/core-ecommerce-system/internal/features/carts"
	"github.com/codepnw/core-ecommerce-system/internal/features/products"
	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/validate"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, userID, addressID string) error
	ListOrders(ctx context.Context, filter *OrderFilter) ([]*OrdersResponse, error)
	UpdateOrderStatus(ctx context.Context, id int64, status OrderStatus) error
}

type OrderServiceConfig struct {
	OrderRepo IOrderRepository          `validate:"required"`
	CartSrv   carts.ICartService        `validate:"required"`
	ProdSrv   products.IProductService  `validate:"required"`
	AddrSrv   addresses.IAddressServide `validate:"required"`
	Tx        *database.TxManager       `validate:"required"`
}

func NewOrderService(cfg *OrderServiceConfig) (IOrderService, error) {
	if err := validate.Struct(cfg); err != nil {
		return nil, fmt.Errorf("OrderServiceConfig required all fields: %w", err)
	}
	return cfg, nil
}

func (s *OrderServiceConfig) CreateOrder(ctx context.Context, userID, addressID string) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	// CART TOTAL PRICE
	var total int64
	products, err := s.CartSrv.GetCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("get cart failed: %w", err)
	}
	if len(products) == 0 {
		return errors.New("cart is empty")
	}
	for _, product := range products {
		total += int64(product.ProductPrice) * product.ProductQuantity
	}

	// GET USER ADDRESS
	addr, err := s.AddrSrv.GetAddressByID(ctx, addressID)
	if err != nil {
		return fmt.Errorf("get address failed: %w", err)
	}

	// TRANSACTION
	err = s.Tx.Transaction(ctx, func(tx *sql.Tx) error {
		// CREATE ORDER
		orderID, err := s.OrderRepo.InsertOrder(ctx, tx, &Order{
			UserID:     userID,
			AddressID:  addr.ID,
			TotalPrice: total,
			Status:     string(StatusPending),
		})
		if err != nil {
			return fmt.Errorf("insert order failed: %w", err)
		}

		// CREATE ORDER ADDRESS
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
			return fmt.Errorf("insert order_address failed: %w", err)
		}

		// CREATE ORDER ITEMS
		var items []*OrderItem
		for _, product := range products {
			ok, err := s.ProdSrv.DeductStock(ctx, tx, product.ProductID, int(product.ProductQuantity))
			if err != nil {
				return fmt.Errorf("deduct product stock failed: %w", err)
			}
			if !ok {
				return fmt.Errorf("product %v out of stock", product.ProductName)
			}

			item := &OrderItem{
				OrderID:   orderID,
				ProductID: product.ProductID,
				Quantity:  int(product.ProductQuantity),
				Price:     int64(product.ProductPrice),
				SubTotal:  int64(product.ProductPrice) * product.ProductQuantity,
			}
			items = append(items, item)
		}
		err = s.OrderRepo.InsertOrderItems(ctx, tx, items)
		if err != nil {
			return fmt.Errorf("insert order_items failed: %w", err)
		}

		// CLEAR CART
		if err = s.CartSrv.ClearCartTx(ctx, tx, userID); err != nil {
			return fmt.Errorf("clear cart failed: %w", err)
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
