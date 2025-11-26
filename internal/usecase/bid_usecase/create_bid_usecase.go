package bid_usecase

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/logger"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/bid_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
)

type BidUsecase struct {
	BidRepository bid_entity.BidRepositoryInterface

	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidChannel          chan bid_entity.Bid
}

type BidUseCaseInterface interface {
	CreateBid(ctx context.Context, bidInputDTO BidInputDTO) *internal_error.InternalError
	FindBidsByAuctionID(ctx context.Context, auctionID string) ([]BidOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionID string) (*BidOutputDTO, *internal_error.InternalError)
}

type BidOutputDTO struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	AuctionId string    `json:"auctionId"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02T15:04:05Z07:00"`
}

type BidInputDTO struct {
	UserId    string  `json:"user_id"`
	AuctionId string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}

var bidBatch []bid_entity.Bid

func NewBidUseCase(bidRepository bid_entity.BidRepositoryInterface) BidUseCaseInterface {
	//Interval and Size from env variables
	maxSizeInterval := getMaxBatchSizeInterval()
	// maxBatchSize defines the maximum number of bids to batch before inserting
	maxBatchSize := getMaxBatchSize()

	bidUseCase := &BidUsecase{
		BidRepository:       bidRepository,
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxSizeInterval,
		timer:               time.NewTimer(maxSizeInterval),
		bidChannel:          make(chan bid_entity.Bid, maxBatchSize),
	}

	// Start the routine to process bid creations
	bidUseCase.triggerCreateRoutine(context.Background())

	return bidUseCase
}

func (bu *BidUsecase) triggerCreateRoutine(ctx context.Context) {
	go func() {
		defer close(bu.bidChannel)

		for {
			select {
			case bidEntity, ok := <-bu.bidChannel:
				if !ok {
					// Channel closed, process remaining bids and exit
					if len(bidBatch) > 0 {
						if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
							logger.Error("error trying to process bid batch list", err)
						}
					}
					return
				}
				// Add bid to the batch
				bidBatch = append(bidBatch, bidEntity)

				// If batch size reached, insert bids
				if len(bidBatch) >= bu.maxBatchSize {
					if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}
					// Reset the batch and timer after processing, for the next batch
					bidBatch = nil
					bu.timer.Reset(bu.batchInsertInterval)
				}
			case <-bu.timer.C:
				// Time interval reached, insert bids if any
				if err := bu.BidRepository.CreateBid(ctx, bidBatch); err != nil {
					logger.Error("error trying to process bid batch list", err)
				}
				bidBatch = nil
				bu.timer.Reset(bu.batchInsertInterval)
			}
		}
	}()
}

func (bu *BidUsecase) CreateBid(ctx context.Context, bidInputDTO BidInputDTO) *internal_error.InternalError {
	bidEntity, err := bid_entity.CreateBid(bidInputDTO.UserId, bidInputDTO.AuctionId, bidInputDTO.Amount)
	if err != nil {
		return err
	}

	bu.bidChannel <- *bidEntity

	return nil
}

func getMaxBatchSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}

func getMaxBatchSize() int {
	value, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 5
	}

	return value
}
