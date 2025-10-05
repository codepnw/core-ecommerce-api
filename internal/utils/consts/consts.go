package consts

import "time"

const ContextTimeout = time.Second * 10

// Params Key
const (
	KeyUserParam    = "user_id"
	KeyAddressParam = "address_id"
)

// Token Exp
const (
	ExpAccessToken  time.Duration = time.Hour * 24
	ExpRefreshToken time.Duration = time.Hour * 24 * 7
)
