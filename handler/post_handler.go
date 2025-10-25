package handler

import (
	"golang-sosmed-gin/dto"
	"golang-sosmed-gin/errorhandler"
	"golang-sosmed-gin/helper"
	"golang-sosmed-gin/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type postHandler struct {
	service service.PostService
}

func NewPostHandler(s service.PostService) *postHandler {
	return &postHandler{
		service: s,
	}
}

func (h *postHandler) Posting(c *gin.Context) {
	var posting dto.PostRequest

	if err := c.ShouldBind(&posting); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: "handler " + err.Error()})
		return
	}

	userID, _ := c.Get("UserID")

	posting.UserID = userID.(int)

	if err := h.service.Posting(&posting); err != nil {
		errorhandler.HandleError(c, &errorhandler.InternalServerError{Message: err.Error()})
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "Posting upload successfully",
	})

	c.JSON(http.StatusCreated, res)
}

func (h *postHandler) MyPost(c *gin.Context) {

	userID, _ := c.Get("UserID")

	idInt := userID.(int)

	post := h.service.MyPost(idInt)

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "My posts",
		Data:       post,
	})

	c.JSON(http.StatusOK, res)
}

func (h *postHandler) DeletePost(c *gin.Context) {

	postID := c.Param("id")
	idInt, err := strconv.Atoi(postID)

	if err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: "handler " + err.Error()})
		return
	}

	erro := h.service.DeletePost(idInt)
	if erro != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: "handler " + erro.Error()})
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Delete successfully",
	})

	c.JSON(http.StatusOK, res)
}

func (h *postHandler) LikePost(c *gin.Context) {

	userID, _ := c.Get("UserID")
	postID := c.PostForm("post_id")

	logrus.Info("post", postID, userID)
	logrus.Info("USERid", userID)
	idInt := userID.(int)
	postInt, _ := strconv.Atoi(postID)

	post := h.service.LikePost(postInt, idInt)
	logrus.Info("post service", post)
	var message string = "Success Like or Unlike"

	var statusCode int = 201

	if post != nil {
		statusCode = http.StatusNotFound
		message = post.Error()
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: statusCode,
		Message:    message,
	})

	c.JSON(statusCode, res)
}
