package main

import "github.com/gin-gonic/gin"

// Router 설정
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/api/:id", getApiHandler)
	return r
}
