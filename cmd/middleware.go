package main

import "github.com/gin-gonic/gin"

// https://zhuanlan.zhihu.com/p/30184285330
func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get session
	}
}
