package server

import (
	"github.com/gin-gonic/gin"
)

func CreateRouter(h *Handler) *gin.Engine {
	r := gin.Default()
	r.POST("/add", h.IncreaseBalance)
	r.POST("/transfer", h.TransferMoney)
	r.GET("/history/:user_id", h.GetUserHistory)

	return r
}
