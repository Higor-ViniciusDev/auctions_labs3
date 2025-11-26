package bid_usecase

import (
	"context"

	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
)

func (bu *BidUsecase) FindBidsByAuctionID(ctx context.Context, auctionID string) ([]BidOutputDTO, *internal_error.InternalError) {
	bidEntities, err := bu.BidRepository.FindBidsByAuctionID(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	var bindOutputDtos []BidOutputDTO
	for _, bidEntity := range bidEntities {
		bidOutputDTO := BidOutputDTO{
			Id:        bidEntity.Id,
			UserId:    bidEntity.UserId,
			AuctionId: bidEntity.AuctionId,
			Amount:    bidEntity.Amount,
			Timestamp: bidEntity.Timestamp,
		}

		bindOutputDtos = append(bindOutputDtos, bidOutputDTO)
	}
	return bindOutputDtos, nil
}

func (bu *BidUsecase) FindWinningBidByAuctionId(ctx context.Context, auctionID string) (*BidOutputDTO, *internal_error.InternalError) {
	winningBidEntity, err := bu.BidRepository.FindWinningBidByAuctionId(ctx, auctionID)

	if err != nil {
		return nil, err
	}

	Bid := &BidOutputDTO{
		Id:        winningBidEntity.Id,
		UserId:    winningBidEntity.UserId,
		AuctionId: winningBidEntity.AuctionId,
		Amount:    winningBidEntity.Amount,
		Timestamp: winningBidEntity.Timestamp,
	}

	return Bid, nil
}
