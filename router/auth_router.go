package router

import (
	"golang-sosmed-gin/config"
	"golang-sosmed-gin/handler"
	"golang-sosmed-gin/repository"
	"golang-sosmed-gin/service"

	"github.com/gin-gonic/gin"
)

func AuthRouter(api *gin.RouterGroup) {
	authRepository := repository.NewAuthRepository(config.DB)
	authService := service.NewAuthService(authRepository)
	authHandler := handler.NewAuthHandler(authService)

	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

}
