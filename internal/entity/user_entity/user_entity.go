package user_entity

import (
	"context"

	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(ctx context.Context, userId string) (*User, *internal_error.InternalError)
}
