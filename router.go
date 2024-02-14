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
	r.GET("/db/tables", getTablesHandler)   // 테이블 리스트 조회
	r.GET("/db/table/:id", getTableHandler) // 테이블 상세정보 조회
	return r
}
