package bid_entity

import (
	"context"
	"time"

	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
	"github.com/google/uuid"
)

type Bid struct {
	Id        string
	UserId    string
	AuctionId string
	Amount    float64
	Timestamp time.Time
}

type BidRepositoryInterface interface {
	FindBidsByAuctionID(ctx context.Context, auctionID string) ([]Bid, *internal_error.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionID string) (*Bid, *internal_error.InternalError)
	CreateBid(ctx context.Context, bidEntities []Bid) *internal_error.InternalError
}

func CreateBid(userID, auctionID string, amount float64) (*Bid, *internal_error.InternalError) {
	bid := &Bid{
		Id:        uuid.New().String(),
		UserId:    userID,
		AuctionId: auctionID,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if err := bid.IsValid(); err != nil {
		return nil, err
	}

	return bid, nil
}

func (b *Bid) IsValid() *internal_error.InternalError {
	if err := uuid.Validate(b.UserId); err != nil {
		return internal_error.NewBadRequestError("UserID is invalid")
	}

	if err := uuid.Validate(b.AuctionId); err != nil {
		return internal_error.NewBadRequestError("AuctionID is invalid")
	}

	if b.Amount <= 0 {
		return internal_error.NewBadRequestError("Amount must be greater than zero")
	}

	return nil
}
