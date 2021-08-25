package customclaims

import "fmt"

type CustomClaimsKey string

var NotFound = fmt.Errorf("not found")