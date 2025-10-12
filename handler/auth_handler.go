package handler

import (
	"golang-sosmed-gin/dto"
	"golang-sosmed-gin/errorhandler"
	"golang-sosmed-gin/helper"
	"golang-sosmed-gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *authHandler {
	return &authHandler{
		service: s,
	}
}

func (h *authHandler) Register(c *gin.Context) {
	var register dto.RegisterRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := h.service.Register(&register); err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "Register successfully, please login",
	})

	c.JSON(http.StatusCreated, res)

}

func (h *authHandler) Login(c *gin.Context) {
	var login dto.LoginRequest

	if err := c.ShouldBindJSON(&login); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{})
		return
	}

	res, err := h.service.Login(&login)

	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	result := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Login successfully",
		Data:       res,
	})

	c.JSON(http.StatusOK, result)

}
