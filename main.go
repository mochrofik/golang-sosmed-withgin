package main

import (
	"fmt"
	"golang-sosmed-gin/config"
	"golang-sosmed-gin/entity"
	"golang-sosmed-gin/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	config.LoadConfig()
	config.ConnectDB()

	const imageDir = "./storage"

	config.DB.AutoMigrate(entity.User{}, entity.Post{}, entity.UploadPosting{}, entity.LikePosting{})

	r := gin.Default()
	r.StaticFS("/storage", http.Dir(imageDir))

	api := r.Group("/api")

	router.UserRouter(api)
	router.AuthRouter(api)
	router.PostRouter(api)

	r.Run(fmt.Sprintf(":%v", config.ENV.PORT)) // listens on 0.0.0.0:8080 by default
}
