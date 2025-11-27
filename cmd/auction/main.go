package main

import (
	"context"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/database/mongodb"
	"github.com/Higor-ViniciusDev/auction_labs3/configuration/logger"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/infra/api/web/controller/auction_controller"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/infra/api/web/controller/bid_controller"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/infra/api/web/controller/user_controller"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/infra/database/auction"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/infra/database/bid"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/infra/database/user"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/usecase/auction_usecase"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/usecase/bid_usecase"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	defer logger.GetLogger().Sync()
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		logger.Error("Erro ao carregar variaveis de ambiente", err)
		return
	}

	db, err := mongodb.NewConnectionDataBaseMongoDB(ctx)
	if err != nil {
		logger.Error(err.Error(), err)
		return
	}
	router := gin.Default()

	userController, auctionController, bidController := initDependeces(db)

	router.GET("/user/:userId", userController.FindUserByID)
	router.GET("/auction", auctionController.FindAuctions)
	router.GET("/auction/:auctionId", auctionController.FindAuctionByID)
	router.POST("/auction", auctionController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionController.FindWinningBidByAuctionId)
	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidsByAuctionID)

	logger.Info("Server iniciando na porta 8080")
	router.Run(":8080")
}

func initDependeces(db *mongo.Database) (user_controller.UserController, auction_controller.AuctionController, bid_controller.BidController) {
	var userController user_controller.UserController
	var auctionController auction_controller.AuctionController
	var bidController bid_controller.BidController

	//User dependeces
	userRepository := user.NewUserRepository(db)
	userUsecase := user_usecase.NewUserUsecase(userRepository)
	userController = *user_controller.NewUserController(userUsecase)

	//Bid dependeces
	bidRepository := bid.NewBidRepository(db)
	bidUsecase := bid_usecase.NewBidUseCase(bidRepository)
	bidController = *bid_controller.NewBidController(bidUsecase)

	//Auction dependeces
	auctionRepository := auction.NewAuctionRepository(db)
	auctionUsecase := auction_usecase.NewAcutionUsecase(auctionRepository, bidRepository)
	auctionController = *auction_controller.NewAuctionController(auctionUsecase)

	return userController, auctionController, bidController
}
