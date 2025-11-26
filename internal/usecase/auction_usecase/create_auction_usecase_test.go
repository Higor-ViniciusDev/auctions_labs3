package auction_usecase

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/auction_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/bid_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
)

type mockAuctionRepo struct {
	createCalled bool
	updateCalled bool
	updateCh     chan struct{}
}

func (m *mockAuctionRepo) CreateAuction(ctx context.Context, auction *auction_entity.Auction) *internal_error.InternalError {
	m.createCalled = true
	return nil
}

func (m *mockAuctionRepo) FindAuctionByID(ctx context.Context, auctionID string) (*auction_entity.Auction, *internal_error.InternalError) {
	return nil, nil
}

func (m *mockAuctionRepo) FindAuctions(ctx context.Context, status auction_entity.ProductCondition, productName, Category string) ([]auction_entity.Auction, *internal_error.InternalError) {
	return nil, nil
}

func (m *mockAuctionRepo) UpdateStatus(ctx context.Context, auctionEntity *auction_entity.Auction, novoStatus auction_entity.AuctionStatus) *internal_error.InternalError {
	m.updateCalled = true
	if m.updateCh != nil {
		select {
		case m.updateCh <- struct{}{}:
		default:
		}
	}
	return nil
}

type mockBidRepo struct{}

func (m *mockBidRepo) FindBidsByAuctionID(ctx context.Context, auctionID string) ([]bid_entity.Bid, *internal_error.InternalError) {
	return nil, nil
}

func (m *mockBidRepo) FindWinningBidByAuctionId(ctx context.Context, auctionID string) (*bid_entity.Bid, *internal_error.InternalError) {
	return nil, nil
}

func (m *mockBidRepo) CreateBid(ctx context.Context, bidEntities []bid_entity.Bid) *internal_error.InternalError {
	return nil
}

func TestCreateAuction_WithGoroutineUpdatesStatus(t *testing.T) {
	// set a short duration for the goroutine timer
	os.Setenv("MAX_INTERVAL_DURATION_AUCTION", "100ms")
	defer os.Unsetenv("MAX_INTERVAL_DURATION_AUCTION")

	mock := &mockAuctionRepo{updateCh: make(chan struct{}, 1)}
	bidMock := &mockBidRepo{}

	usecase := NewAcutionUsecase(mock, bidMock)

	ctx := context.Background()

	input := &AuctionInputDTO{
		ProductName: "produtoteste",
		Category:    "teste",
		Description: "teste do teste",
		Condition:   ProductCondition(auction_entity.New),
	}

	if err := usecase.CreateAuction(ctx, input); err != nil {
		t.Fatalf("CreateAuction returned error: %v", err)
	}

	// wait for goroutine to call UpdateStatus
	select {
	case <-mock.updateCh:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected UpdateStatus to be called by goroutine")
	}

	if !mock.createCalled {
		t.Fatal("expected CreateAuction to call repository CreateAuction")
	}
}
