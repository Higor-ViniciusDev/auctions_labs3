package bid

import (
	"context"
	"fmt"
	"sync"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/logger"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/auction_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/bid_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/infra/database/auction"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BidEntityMongo struct {
	ID        string  `bson:"_id"`
	AuctionID string  `bson:"auction_id"`
	UserID    string  `bson:"user_id"`
	Amount    float64 `bson:"amount"`
	Timestamp int64   `bson:"timestamp"`
}

type BidRepository struct {
	Collection  *mongo.Collection
	AuctionRepo *auction.AuctionRepository
}

func NewBidRepository(db *mongo.Database) *BidRepository {
	return &BidRepository{
		Collection: db.Collection("bids"),
	}
}

func (br *BidRepository) CreateBid(ctx context.Context, bidEntities []bid_entity.Bid) *internal_error.InternalError {
	var wg sync.WaitGroup

	for _, bid := range bidEntities {
		wg.Add(1)

		go func(bid bid_entity.Bid) {
			defer wg.Done()

			auctionEntity, err := br.AuctionRepo.FindAuctionByID(ctx, bid.AuctionId)

			if err != nil {
				logger.Error(fmt.Sprintf("Failed to find auction for bid ID %s", bid.Id), err)
				//Ignorar ação pois não achou o leilão
				return
			}

			if auctionEntity.Status != auction_entity.Active {
				logger.Info(fmt.Sprintf("Leilão do id: %v foi feito lances porém ele não esta ativo", bid.AuctionId))
				//Ignorar ação pois o leilão não está ativo
				return
			}

			bidEntityMongo := BidEntityMongo{
				ID:        bid.Id,
				AuctionID: bid.AuctionId,
				UserID:    bid.UserId,
				Amount:    bid.Amount,
				Timestamp: bid.Timestamp.Unix(),
			}
			_, errorInsert := br.Collection.InsertOne(ctx, bidEntityMongo)

			if errorInsert != nil {
				dataParam := map[string]interface{}{
					"bid_id":     bid.Id,
					"auction_id": bid.AuctionId,
				}

				logger.Error(fmt.Sprintf("Failed to insert bid\n param: %v", dataParam), errorInsert)
			}
		}(bid)
	}
	wg.Wait()
	return nil
}
