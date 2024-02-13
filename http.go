package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func StartServer() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 환경변수에서 포트 가져오기
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // 기본포트 설정
	}

	r := SetupRouter()
	r.Run(":" + port)

}
