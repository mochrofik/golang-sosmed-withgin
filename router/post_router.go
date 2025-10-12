package router

import (
	"golang-sosmed-gin/config"
	"golang-sosmed-gin/handler"
	"golang-sosmed-gin/middleware"
	"golang-sosmed-gin/repository"
	"golang-sosmed-gin/service"

	"github.com/gin-gonic/gin"
)

func PostRouter(api *gin.RouterGroup) {
	postRepository := repository.NewPostRepository(config.DB)

	postService := service.NewPostService(postRepository)

	postHandler := handler.NewPostHandler(postService)

	r := api.Group("/posting")

	r.Use(middleware.JWTMiddleware())

	r.POST("/", postHandler.Posting)
	r.GET("/my-post", postHandler.MyPost)
}
