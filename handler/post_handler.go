package handler

import (
	"golang-sosmed-gin/dto"
	"golang-sosmed-gin/errorhandler"
	"golang-sosmed-gin/helper"
	"golang-sosmed-gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
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
		StatusCode: http.StatusCreated,
		Message:    "My posts",
		Data:       post,
	})

	c.JSON(http.StatusCreated, res)

}
