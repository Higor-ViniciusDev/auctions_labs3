package auction_controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/rest_err"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/usecase/auction_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ac *AuctionController) FindAuctionByID(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		restError := rest_err.NewBadRequestError("auctionId invalid")
		c.JSON(restError.Code, restError.Message)
		return
	}

	auction, err := ac.AuctionUsecase.FindAuctionByID(c.Request.Context(), auctionId)
	if err != nil {
		restErrorConvertido := rest_err.ConvertInternalErrorToRestError(err)
		c.JSON(restErrorConvertido.Code, restErrorConvertido.Message)
		return
	}

	c.JSON(http.StatusOK, auction)
}

func (ac *AuctionController) FindAuctions(c *gin.Context) {
	status := c.Query("status")
	nameProduct := c.Query("productName")
	category := c.Query("category")

	statusNumber, errConv := strconv.Atoi(status)

	if errConv != nil {
		restError := rest_err.NewBadRequestError("error convert string for number, param invalid: status")
		c.JSON(restError.Code, restError.Message)
		return
	}

	users, err := ac.AuctionUsecase.FindAuctions(c.Request.Context(), auction_usecase.ProductCondition(statusNumber), nameProduct, category)
	if err != nil {
		restErrorConvertido := rest_err.ConvertInternalErrorToRestError(err)
		c.JSON(restErrorConvertido.Code, restErrorConvertido.Message)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (u *AuctionController) FindWinningBidByAuctionId(c *gin.Context) {
	auctionId := c.Param("auctionId")

	if err := uuid.Validate(auctionId); err != nil {
		errRest := rest_err.NewBadRequestError("Invalid fields", rest_err.Causes{
			Field:   "auctionId",
			Message: "Invalid UUID value",
		})

		c.JSON(errRest.Code, errRest)
		return
	}

	auctionData, err := u.AuctionUsecase.FindWinningBidByAuctionId(context.Background(), auctionId)
	if err != nil {
		errRest := rest_err.ConvertInternalErrorToRestError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, auctionData)
}
