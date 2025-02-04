package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/dayterr/test-go-iq/internal/storage"
)

type Handler struct {
	Storage storage.Storage
}

func NewHandler() *Handler {
	h := Handler{}
	return &h
}

func (h *Handler) IncreaseBalance(c *gin.Context) {
	var u storage.User

	err := c.ShouldBindBodyWithJSON(&u)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if u.ID == 0 || u.Balance == 0 {
		c.String(http.StatusBadRequest, "a field is missing")
		return
	}

	err = h.Storage.IncreaseBalance(u)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, "Success")

}

func (h *Handler) TransferMoney(c *gin.Context) {
	var t storage.Transaction

	err := c.ShouldBindBodyWithJSON(&t)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if t.SenderID == 0 || t.Amount == 0 || t.RecieverID == 0 {
		c.String(http.StatusBadRequest, "a field is missing")
		return
	}

	err = h.Storage.TransferMoney(t)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusCreated, "Success")

}

func (h *Handler) GetUserHistory(c *gin.Context) {
	var userID int

	userID, err := strconv.Atoi(c.Params.ByName("user_id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if userID <= 0 {
		c.String(http.StatusBadRequest, "user id must be positive")
		return
	}

	operations, err := h.Storage.GetHistory(userID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, operations)
}
