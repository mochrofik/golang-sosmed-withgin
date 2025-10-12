package router

import (
	"golang-sosmed-gin/config"
	"golang-sosmed-gin/handler"
	"golang-sosmed-gin/middleware"
	"golang-sosmed-gin/repository"
	"golang-sosmed-gin/service"

	"github.com/gin-gonic/gin"
)

func UserRouter(api *gin.RouterGroup) {

	userRepository := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepository)
	authHandler := handler.NewUserHandler(userService)

	r := api.Group("/user")

	r.Use(middleware.JWTMiddleware())

	r.GET("/all", authHandler.GetAllUser)

}
