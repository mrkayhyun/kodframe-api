package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

// Router 설정
func SetupRouter() *gin.Engine {
	// 데이터베이스 연결 초기화
	if err := InitDB(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := gin.Default()
	// api
	r.GET("/api/:id", getApiHandler)
	// db
	r.GET("/db/tables", getTablesHandler)
	r.GET("/db/table/:id", getTableHandler)
	return r
}
