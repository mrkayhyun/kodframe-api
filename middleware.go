package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

// 접근을 허용하고 싶은 IP 가 있을경우 설정한다.
// 설정 하기 않을경우 모든 IP를 허용 하고 설정 할 경우 설정된 IP만 접근을 허용한다.
func IpAccessMiddleware(c *gin.Context) {
	// Todo : 환경변수에서 허용 IP를 설정 하게 되어 있는데 DB로 변경하는건 추후 고려
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ipPermitStr := os.Getenv("IP_PERMIT")
	if len(ipPermitStr) > 0 {
		// 설정된 IP만 허용
		ipPermit := strings.Split(ipPermitStr, ",")
		if !utils.Contains(ipPermit, c.RemoteIP()) {
			c.AbortWithStatus(http.StatusForbidden) // 설정 되지 않은 아이피 차단
			return
		}
	}
	c.Next()
}
