package auction

import (
	"context"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/logger"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/auction_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (a *AuctionRepository) UpdateStatus(ctx context.Context, auctionEntity *auction_entity.Auction, novoStatus auction_entity.AuctionStatus) *internal_error.InternalError {
	filter := bson.D{{"_id", auctionEntity.Id}}
	update := bson.D{{"$set", bson.D{{"status", novoStatus}}}}

	_, err := a.Collection.UpdateOne(ctx, filter, update)

	if err != nil {
		logger.Error("failed to update auction in repository", err)
		return internal_error.NewInternalServerError("failed to update auction: " + err.Error())
	}

	return nil
}
