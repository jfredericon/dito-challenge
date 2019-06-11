package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Home : GET "/" endpoint
func Home() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Dito Challenge API",
		})

	}

}
