package bid_controller

import (
	"net/http"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/rest_err"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/infra/api/web/validation"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/usecase/bid_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BidController struct {
	BidUsecase bid_usecase.BidUseCaseInterface
}

func NewBidController(BidUsecase bid_usecase.BidUseCaseInterface) *BidController {
	return &BidController{
		BidUsecase: BidUsecase,
	}
}

func (bc *BidController) CreateBid(c *gin.Context) {
	var bidInput bid_usecase.BidInputDTO

	if err := c.ShouldBindJSON(&bidInput); err != nil {
		restErro := validation.ValidateErr(err)

		c.JSON(restErro.Code, restErro.Message)
		return
	}

	err := bc.BidUsecase.CreateBid(c.Request.Context(), bidInput)

	if err != nil {
		restErro := rest_err.ConvertInternalErrorToRestError(err)

		c.JSON(restErro.Code, restErro.Message)
		return
	}

	c.Status(http.StatusCreated)
}

func (bc *BidController) FindBidsByAuctionID(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		restError := rest_err.NewBadRequestError("auctionId invalid")
		c.JSON(restError.Code, restError.Message)
		return
	}

	auctions, err := bc.BidUsecase.FindBidsByAuctionID(c.Request.Context(), auctionId)
	if err != nil {
		restErrorConvertido := rest_err.ConvertInternalErrorToRestError(err)
		c.JSON(restErrorConvertido.Code, restErrorConvertido.Message)
		return
	}

	c.JSON(http.StatusOK, auctions)
}

func (bc *BidController) FindWinningBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		restError := rest_err.NewBadRequestError("auctionId invalid")
		c.JSON(restError.Code, restError.Message)
		return
	}

	auction, err := bc.BidUsecase.FindWinningBidByAuctionId(c.Request.Context(), auctionId)
	if err != nil {
		restErrorConvertido := rest_err.ConvertInternalErrorToRestError(err)
		c.JSON(restErrorConvertido.Code, restErrorConvertido.Message)
		return
	}

	c.JSON(http.StatusOK, auction)
}
