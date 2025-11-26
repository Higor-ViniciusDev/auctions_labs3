package auction_controller

import (
	"net/http"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/rest_err"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/infra/api/web/validation"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/usecase/auction_usecase"
	"github.com/gin-gonic/gin"
)

type AuctionController struct {
	AuctionUsecase auction_usecase.AuctionUsecaseInterface
}

func NewAuctionController(auctionUsecase auction_usecase.AuctionUsecaseInterface) *AuctionController {
	return &AuctionController{
		AuctionUsecase: auctionUsecase,
	}
}

func (ac *AuctionController) CreateAuction(c *gin.Context) {
	var auctioninput auction_usecase.AuctionInputDTO

	if err := c.ShouldBindJSON(&auctioninput); err != nil {
		restErro := validation.ValidateErr(err)

		c.JSON(restErro.Code, restErro.Message)
		return
	}

	err := ac.AuctionUsecase.CreateAuction(c.Request.Context(), &auctioninput)

	if err != nil {
		restErro := rest_err.ConvertInternalErrorToRestError(err)

		c.JSON(restErro.Code, restErro.Message)
		return
	}

	c.Status(http.StatusCreated)
}
