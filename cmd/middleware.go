package main

import "github.com/gin-gonic/gin"

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get session
	}
}
