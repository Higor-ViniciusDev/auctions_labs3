package auction_usecase

import (
	"context"
	"os"
	"time"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/logger"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/auction_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/bid_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/usecase/bid_usecase"
)

type AuctionUsecase struct {
	AuctionRepostiry auction_entity.AuctionRepositoryInterface
	BidRepository    bid_entity.BidRepositoryInterface
}

type AuctionUsecaseInterface interface {
	CreateAuction(ctx context.Context, auction *AuctionInputDTO) *internal_error.InternalError
	FindAuctionByID(ctx context.Context, auctionID string) (*AuctionOutputDTO, *internal_error.InternalError)
	FindAuctions(ctx context.Context, status ProductCondition, productName, Category string) ([]AuctionOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionID string) (*WinningInfoOutputDTO, *internal_error.InternalError)
}

type AuctionOutputDTO struct {
	Id          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	TimeStamp   time.Time        `json:"timestamp" time_format:"2006-01-02T15:04:05Z07:00"`
}

type AuctionInputDTO struct {
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
}

type ProductCondition int64
type AuctionStatus int64

type WinningInfoOutputDTO struct {
	AuctionId string                    `json:"auctionId"`
	Bid       *bid_usecase.BidOutputDTO `json:"bid,omitempty"`
}

func NewAcutionUsecase(auctionRepository auction_entity.AuctionRepositoryInterface, bidRepository bid_entity.BidRepositoryInterface) *AuctionUsecase {
	useCaseAuction := &AuctionUsecase{
		AuctionRepostiry: auctionRepository,
		BidRepository:    bidRepository,
	}

	// go useCaseAuction.TriggerFindExpiredAuctions(context.Background())
	return useCaseAuction
}

func (au *AuctionUsecase) CreateAuction(ctx context.Context, auction *AuctionInputDTO) *internal_error.InternalError {
	auctionEntity, err := auction_entity.NewAuctionEntity(
		auction.ProductName,
		auction.Category,
		auction.Description,
		auction_entity.ProductCondition(auction.Condition),
	)

	if err != nil {
		return err
	}

	errorCreate := au.AuctionRepostiry.CreateAuction(ctx, auctionEntity)

	if errorCreate != nil {
		return errorCreate
	}

	go func(auction *auction_entity.Auction) {
		//Interval from env variables
		maxSizeInterval := getMaxDurationAuction()
		timer := time.NewTimer(maxSizeInterval)

		<-timer.C
		err := au.AuctionRepostiry.UpdateStatus(ctx, auction, auction_entity.Completed)
		if err != nil {
			logger.Error("error trying to update status auction", err)
		}
	}(auctionEntity)

	return nil
}

func getMaxDurationAuction() time.Duration {
	batchInsertInterval := os.Getenv("MAX_INTERVAL_DURATION_AUCTION")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 30 * time.Second
	}

	return duration
}
