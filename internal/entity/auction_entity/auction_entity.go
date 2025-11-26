package auction_entity

import (
	"context"
	"time"

	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
	"github.com/google/uuid"
)

type Auction struct {
	Id          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
	TimeStamp   time.Time
}

type ProductCondition int
type AuctionStatus int

const (
	Active AuctionStatus = iota
	Completed
)

const (
	New ProductCondition = iota
	Used
	Refurbished
)

func (a *Auction) IsValid() *internal_error.InternalError {
	if (len(a.ProductName) <= 3 || len(a.Category) <= 2 || len(a.Description) <= 5) && (a.Condition != New || a.Condition != Used || a.Condition != Refurbished) {
		return internal_error.NewBadRequestError("invalid auction entity")
	}

	return nil
}

func NewAuctionEntity(productName, category, description string, condition ProductCondition) (*Auction, *internal_error.InternalError) {
	auction := &Auction{
		Id:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   ProductCondition(condition),
		Status:      AuctionStatus(New),
		TimeStamp:   time.Now(),
	}

	if err := auction.IsValid(); err != nil {
		return nil, err
	}

	return auction, nil
}

type AuctionRepositoryInterface interface {
	CreateAuction(ctx context.Context, auction *Auction) *internal_error.InternalError
	FindAuctionByID(ctx context.Context, auctionID string) (*Auction, *internal_error.InternalError)
	FindAuctions(ctx context.Context, status ProductCondition, productName, Category string) ([]Auction, *internal_error.InternalError)
	UpdateStatus(ctx context.Context, auctionEntity *Auction, novoStatus AuctionStatus) *internal_error.InternalError
}
