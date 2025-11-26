package user_usecase

import (
	"context"

	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/user_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
)

type UserUsecase struct {
	UserRepository user_entity.UserRepositoryInterface
}

func NewUserUsecase(user_repository user_entity.UserRepositoryInterface) *UserUsecase {
	return &UserUsecase{
		UserRepository: user_repository,
	}
}

type UserOutputDTO struct {
	Id   string
	Name string
}

type UserUsecaseInterface interface {
	FindUserById(ctx context.Context, userId string) (*UserOutputDTO, *internal_error.InternalError)
}

func (u *UserUsecase) FindUserById(ctx context.Context, userID string) (*UserOutputDTO, *internal_error.InternalError) {
	user, err := u.UserRepository.FindUserById(ctx, userID)

	if err != nil {
		return nil, err
	}

	return &UserOutputDTO{
		Id:   user.Id,
		Name: user.Name,
	}, nil
}
