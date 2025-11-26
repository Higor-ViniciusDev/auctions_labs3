package auction

import (
	"context"
	"fmt"
	"time"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/logger"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/auction_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (ar *AuctionRepository) FindAuctionByID(ctx context.Context, auctionID string) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": auctionID}

	var auctionEntityMongo AuctionEntityMongo

	err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			logger.Error(fmt.Sprintf("auction not found with id: %s", auctionID), err)

			return nil, internal_error.NewNotFoundError("auction not found with id: " + auctionID)
		}

		logger.Error(fmt.Sprintf("internal server failed to find auction by id repository: %s", auctionID), err)
		return nil, internal_error.NewInternalServerError("failed to find auction: " + err.Error())
	}

	auction := &auction_entity.Auction{
		Id:          auctionEntityMongo.Id,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		TimeStamp:   time.Unix(auctionEntityMongo.TimeStamp, 0),
	}

	return auction, nil
}

func (ar *AuctionRepository) FindAuctions(ctx context.Context, status auction_entity.ProductCondition, productName, Category string) ([]auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	//Filter com regex to make case insensitive search
	if productName != "" {
		filter["product_name"] = primitive.Regex{Pattern: productName, Options: "i"}
	}

	if Category != "" {
		filter["category"] = primitive.Regex{Pattern: Category, Options: "i"}
	}

	// Find auctions based on the constructed filter
	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("internal server failed to find auctions in repository", err)
		return nil, internal_error.NewInternalServerError("failed to find auctions: " + err.Error())
	}

	var auctionsMongo []AuctionEntityMongo
	if err = cursor.All(ctx, &auctionsMongo); err != nil {
		logger.Error("internal server failed to decode auctions in repository", err)
		return nil, internal_error.NewInternalServerError("failed to decode auctions: " + err.Error())
	}

	// Map MongoDB entities to domain entities
	var auctions []auction_entity.Auction
	for _, auctionMongo := range auctionsMongo {
		auction := auction_entity.Auction{
			Id:          auctionMongo.Id,
			ProductName: auctionMongo.ProductName,
			Category:    auctionMongo.Category,
			Description: auctionMongo.Description,
			Condition:   auctionMongo.Condition,
			Status:      auctionMongo.Status,
			TimeStamp:   time.Unix(auctionMongo.TimeStamp, 0),
		}
		// Append to the result slice
		auctions = append(auctions, auction)
	}

	return auctions, nil
}
