package user

import (
	uiptypes "github.com/HyperService-Consortium/go-uip/types"
	types "github.com/HyperService-Consortium/go-ves/types"
)

var testuser_impl types.User = &User{}
var testuser_impl2 types.User = User{}

var testaccount_impl uiptypes.Account = XORMUserAdapter{}
var testaccount_impl2 uiptypes.Account = &XORMUserAdapter{}

var testbase_impl types.UserBase = &XORMUserBase{}
var testbase_impl2 types.UserBase = XORMUserBase{}
