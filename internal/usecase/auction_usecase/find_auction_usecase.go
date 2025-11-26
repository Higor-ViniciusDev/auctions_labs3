package auction_usecase

import (
	"context"

	"github.com/Higor-ViniciusDev/auction_labs3/internal/entity/auction_entity"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/internal_error"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/usecase/bid_usecase"
)

func (au *AuctionUsecase) FindAuctionByID(ctx context.Context, auctionID string) (*AuctionOutputDTO, *internal_error.InternalError) {
	auctionEntity, err := au.AuctionRepostiry.FindAuctionByID(ctx, auctionID)

	if err != nil {
		return nil, err
	}

	auctionOutput := &AuctionOutputDTO{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   ProductCondition(auctionEntity.Condition),
		Status:      AuctionStatus(auctionEntity.Status),
		TimeStamp:   auctionEntity.TimeStamp,
	}

	return auctionOutput, nil
}

func (au *AuctionUsecase) FindAuctions(ctx context.Context, status ProductCondition, productName, Category string) ([]AuctionOutputDTO, *internal_error.InternalError) {
	auctionEntities, err := au.AuctionRepostiry.FindAuctions(ctx, auction_entity.ProductCondition(status), productName, Category)
	if err != nil {
		return nil, err
	}

	var auctionsOutput []AuctionOutputDTO
	for _, auctionEntity := range auctionEntities {
		//Criando o DTO de saída
		auctionOutput := AuctionOutputDTO{
			Id:          auctionEntity.Id,
			ProductName: auctionEntity.ProductName,
			Category:    auctionEntity.Category,
			Description: auctionEntity.Description,
			Condition:   ProductCondition(auctionEntity.Condition),
			Status:      AuctionStatus(auctionEntity.Status),
			TimeStamp:   auctionEntity.TimeStamp,
		}

		//Adicionando o DTO de saída na lista de saídas
		auctionsOutput = append(auctionsOutput, auctionOutput)
	}

	return auctionsOutput, nil
}

func (bu *AuctionUsecase) FindWinningBidByAuctionId(ctx context.Context, auctionID string) (*WinningInfoOutputDTO, *internal_error.InternalError) {
	winningBidEntity, err := bu.BidRepository.FindWinningBidByAuctionId(ctx, auctionID)

	if err != nil {
		if err.Err == internal_error.NotFoundError {
			return &WinningInfoOutputDTO{
				AuctionId: auctionID,
				Bid:       nil}, nil
		}

		return nil, err
	}

	winningInfoOutputDTO := &WinningInfoOutputDTO{
		AuctionId: auctionID,
		Bid: &bid_usecase.BidOutputDTO{
			Id:        winningBidEntity.Id,
			UserId:    winningBidEntity.UserId,
			AuctionId: winningBidEntity.AuctionId,
			Amount:    winningBidEntity.Amount,
			Timestamp: winningBidEntity.Timestamp,
		},
	}

	return winningInfoOutputDTO, nil
}
