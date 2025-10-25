package handler

import (
	"golang-sosmed-gin/dto"
	"golang-sosmed-gin/errorhandler"
	"golang-sosmed-gin/helper"
	"golang-sosmed-gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *userHandler {
	return &userHandler{
		service: s,
	}
}

func (h *userHandler) GetAllUser(c *gin.Context) {
	var user dto.UserRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: "handler" + err.Error()})
		return
	}

	res := h.service.GetAllUser(&user)

	usersResponse := make([]dto.UserResponse, 0, len(*res))

	for _, user := range *res {
		// Panggil fungsi helper untuk konversi
		responseItem := dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
		// Input/Tambahkan ke slice tujuan
		usersResponse = append(usersResponse, responseItem)
	}

	result := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Data user",
		Data:       usersResponse,
		Paginate: &dto.Paginate{
			Page:      1,
			PerPage:   10,
			TotalPage: 1,
			Total:     1,
		},
	})

	c.JSON(http.StatusOK, result)
}

func (h *userHandler) MyProfile(c *gin.Context) {

	userID, _ := c.Get("UserID")

	idInt := userID.(int)

	data := h.service.GetMyProfile(idInt)

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "My Profile",
		Data:       data,
	})

	c.JSON(http.StatusOK, res)
}

func (h *userHandler) EditProfile(c *gin.Context) {

	var user dto.ProfileRequest

	userID, _ := c.Get("UserID")

	user.UserID = userID.(int)

	if err := c.ShouldBind(&user); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: "handler " + err.Error()})
		return
	}

	err := h.service.EditProfile(&user)
	var res any

	if err != nil {
		res = helper.Response(dto.ResponseParams{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		c.JSON(http.StatusBadRequest, res)

	} else {
		res = helper.Response(dto.ResponseParams{
			StatusCode: http.StatusCreated,
			Message:    "Profile update successfully",
		})
		c.JSON(http.StatusCreated, res)
	}

}
