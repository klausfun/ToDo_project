package handler

import (
	"github.com/gin-gonic/gin"
	todo "github.com/klausfun/ToDo_project"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		// StatusBadRequest == 400 - означает, что пользователь предоставил некорректные данные в запросе
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	// передаем данные "ниже"(в сервис, где реализована бизнес логика регистрации)
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		// StatusInternalServerError == 500 - обозначает клиенту про внутреннюю ошибку на сервере
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		// StatusBadRequest == 400 - означает, что пользователь предоставил некорректные данные в запросе
		newErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	// передаем данные "ниже"(в сервис, где реализована бизнес логика регистрации)
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		// StatusInternalServerError == 500 - обозначает клиенту про внутреннюю ошибку на сервере
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"token": token})
}
