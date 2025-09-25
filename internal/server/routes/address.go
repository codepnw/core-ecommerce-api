package routes

import (
	"fmt"

	"github.com/codepnw/core-ecommerce-system/internal/features/addresses"
	"github.com/codepnw/core-ecommerce-system/internal/utils/consts"
)

func (cfg *RoutesConfig) addressRoutes() {
	// Setup
	repo := addresses.NewAddressRepository(cfg.DB)
	service := addresses.NewAddressSerivce(repo)
	handler := addresses.NewAddressHandler(service)

	// Params Key
	addressParamKey := fmt.Sprintf("/:%s", consts.KeyAddressParam)
	userParamKey := fmt.Sprintf("/:%s", consts.KeyUserParam)

	r := cfg.Router.Group(cfg.Prefix + "/addresses")

	r.Post("/", handler.CreateAddress)
	r.Get(userParamKey, handler.GetAddressByUserID)
	r.Get(addressParamKey, handler.GetAddressByID)
	r.Patch(addressParamKey, handler.UpdateAddress)
	r.Delete(addressParamKey, handler.DeleteAddress)
	r.Delete(addressParamKey+"/default", handler.SetAddressDefault)
}
