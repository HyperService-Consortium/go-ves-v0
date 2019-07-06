package user

import types "github.com/Myriad-Dreamin/go-ves/types"

var testuser_impl types.User = &User{}
var testuser_impl2 types.User = User{}

var testaccount_impl types.Account = XORMUserAdapter{}
var testaccount_impl2 types.Account = &XORMUserAdapter{}

var testbase_impl types.UserBase = &UserBase{}
var testbase_impl2 types.UserBase = UserBase{}
