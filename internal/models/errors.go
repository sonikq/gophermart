package models

import "errors"

var (
	ErrNotAuthenticated                  = errors.New("user is not authenticated")
	ErrInvalidCredentials                = errors.New("invalid login or password")
	ErrEmptyLogin                        = errors.New("empty login is not allowed")
	ErrEmptyPassword                     = errors.New("empty password is not allowed")
	ErrLoginIsAlreadyTaken               = errors.New("login is already taken")
	ErrNotUniqueOrderNum                 = errors.New("order already exists")
	ErrInvalidOrderNum                   = errors.New("invalid order number format")
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("order has already been uploaded by another user")
	ErrOrderAlreadyUploadedByThisUser    = errors.New("order has already been uploaded by this user")
	ErrEmptyOrderList                    = errors.New("empty order list")
	ErrInsufficientFunds                 = errors.New("insufficient funds in the account")
)
