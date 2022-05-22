package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/model"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var request model.User

	if err := c.BindJSON(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.AuthService.CreateUser(request)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

type signInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var request signInRequest

	if err := c.BindJSON(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.AuthService.GenerateToken(request.Username, request.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
