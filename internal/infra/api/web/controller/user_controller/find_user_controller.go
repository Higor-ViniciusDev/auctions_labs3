package user_controller

import (
	"net/http"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/rest_err"
	"github.com/Higor-ViniciusDev/auction_labs3/internal/usecase/user_usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	UseCaseUser user_usecase.UserUsecaseInterface
}

func NewUserController(useCaseUser user_usecase.UserUsecaseInterface) *UserController {
	return &UserController{
		UseCaseUser: useCaseUser,
	}
}

func (uc *UserController) FindUserByID(c *gin.Context) {
	userId := c.Param("userId")

	if err := uuid.Validate(userId); err != nil {
		restError := rest_err.NewBadRequestError("userId inválido")
		c.JSON(restError.Code, restError.Message)
		return
	}

	user, err := uc.UseCaseUser.FindUserById(c.Request.Context(), userId)
	if err != nil {
		restErrorConvertido := rest_err.ConvertInternalErrorToRestError(err)
		c.JSON(restErrorConvertido.Code, restErrorConvertido.Message)
		return
	}

	c.JSON(http.StatusOK, user)
}
