package storage_handler

import "github.com/HyperService-Consortium/go-ves/types"

var _ types.StorageHandler = new(Database)
