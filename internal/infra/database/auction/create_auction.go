package auction

import (
	"context"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/logger"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/auction_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	TimeStamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(db *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: db.Collection("auctions"),
	}
}

func (a *AuctionRepository) CreateAuction(ctx context.Context, auction *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		TimeStamp:   auction.TimeStamp.Unix(),
	}
	_, err := a.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("failed to create auction in repository", err)
		return internal_error.NewInternalServerError("failed to create auction: " + err.Error())
	}
	return nil
}
