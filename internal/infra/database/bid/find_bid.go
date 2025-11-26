package bid

import (
	"context"
	"fmt"
	"time"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/logger"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/bid_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (br *BidRepository) FindBidsByAuctionID(ctx context.Context, auctionID string) ([]bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionID}

	cursor, err := br.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to find bids for auction ID %s", auctionID), err)
		return nil, internal_error.NewInternalServerError("Failed to find bids")
	}
	defer cursor.Close(ctx)

	var bids []bid_entity.Bid
	for cursor.Next(ctx) {
		var bidEntityMongo BidEntityMongo
		if err := cursor.Decode(&bidEntityMongo); err != nil {
			logger.Error("Failed to decode bid entity", err)
			return nil, internal_error.NewInternalServerError("Failed to decode bid entity")
		}
		bid := bid_entity.Bid{
			Id:        bidEntityMongo.ID,
			UserId:    bidEntityMongo.UserID,
			AuctionId: bidEntityMongo.AuctionID,
			Amount:    bidEntityMongo.Amount,
			Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
		}
		bids = append(bids, bid)
	}

	return bids, nil
}

func (br *BidRepository) FindWinningBidByAuctionId(ctx context.Context, auctionID string) (*bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionID}
	options := options.FindOne().SetSort(bson.D{{Key: "amount", Value: -1}})

	var bidEntityMongo BidEntityMongo
	err := br.Collection.FindOne(ctx, filter, options).Decode(&bidEntityMongo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, internal_error.NewNotFoundError("No bids found for the auction")
		}

		logger.Error(fmt.Sprintf("Failed to find winning bid for auction ID %s", auctionID), err)
		return nil, internal_error.NewInternalServerError("Failed to find winning bid")
	}

	bidEntity := &bid_entity.Bid{
		Id:        bidEntityMongo.ID,
		UserId:    bidEntityMongo.UserID,
		AuctionId: bidEntityMongo.AuctionID,
		Amount:    bidEntityMongo.Amount,
		Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
	}
	return bidEntity, nil
}
