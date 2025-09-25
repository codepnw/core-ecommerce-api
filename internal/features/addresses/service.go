package addresses

import (
	"context"
	"errors"

	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
	"github.com/codepnw/core-ecommerce-system/internal/utils/errs"
	"github.com/gofiber/fiber/v2/log"
)

type IAddressServide interface {
	CreateAddress(ctx context.Context, req *AddressCreate) error
	GetAddressByID(ctx context.Context, id int64) (*Address, error)
	GetAddressByUserID(ctx context.Context, userID string) ([]*Address, error)
	UpdateAddress(ctx context.Context, id int64, req *AddressUpdate) error
	DeleteAddress(ctx context.Context, id int64) error
}

type addressService struct {
	repo IAddressRepository
}

func NewAddressSerivce(repo IAddressRepository) IAddressServide {
	return &addressService{repo: repo}
}

func (s *addressService) CreateAddress(ctx context.Context, req *AddressCreate) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	add := &Address{
		UserID:      req.UserID,
		AddressLine: req.AddressLine,
		City:        req.City,
		State:       req.State,
		PostalCode:  req.PostalCode,
		Phone:       req.Phone,
		IsDefault:   false,
	}

	if err := s.repo.Create(ctx, add); err != nil {
		log.Errorf("create address failed: %v", err)
		return errors.New("create address failed")
	}

	return nil
}

func (s *addressService) GetAddressByID(ctx context.Context, id int64) (*Address, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	res, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrAddressNotFound) {
			return nil, err
		}
		log.Errorf("get address failed: %v", err)
		return nil, errors.New("get address failed")
	}

	return res, nil
}

func (s *addressService) GetAddressByUserID(ctx context.Context, userID string) ([]*Address, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	res, err := s.repo.List(ctx, userID)
	if err != nil {
		log.Errorf("get addresses failed: %v", err)
		return nil, errors.New("get addresses failed")
	}

	return res, nil
}

func (s *addressService) UpdateAddress(ctx context.Context, id int64, req *AddressUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	if err := s.repo.Update(ctx, id, req); err != nil {
		if errors.Is(err, errs.ErrAddressNotFound) {
			return err
		}
		log.Errorf("update address failed: %v", err)
		return errors.New("update address failed")
	}

	return nil
}

func (s *addressService) DeleteAddress(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, consts.ContextTimeout)
	defer cancel()

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, errs.ErrAddressNotFound) {
			return err
		}
		log.Errorf("delete address failed: %v", err)
		return errors.New("delete address failed")
	}

	return nil
}
