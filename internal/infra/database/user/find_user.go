package user

import (
	"context"
	"fmt"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/logger"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/user_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserEntityMongo struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: db.Collection("users"),
	}
}

func (u *UserRepository) FindUserById(ctx context.Context, userId string) (*user_entity.User, *internal_error.InternalError) {
	filter := bson.M{"_id": userId}

	var userEntityMongo UserEntityMongo

	err := u.Collection.FindOne(ctx, filter).Decode(&userEntityMongo)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			logger.Error(fmt.Sprintf("user not found with id: %s", userId), err)

			return nil, internal_error.NewNotFoundError("user not found with id: " + userId)
		}

		logger.Error(fmt.Sprintf("internal server failed to find user by id repository: %s", userId), err)
		return nil, internal_error.NewInternalServerError("failed to find user: " + err.Error())
	}

	user := &user_entity.User{
		Id:   userEntityMongo.Id,
		Name: userEntityMongo.Name,
	}

	return user, nil
}
