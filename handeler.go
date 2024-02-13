package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func getApiHandler(g *gin.Context) {
	id := g.Param("id")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 환경변수에서 포트 가져오기
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// MySQL 데이터베이스 연결 설정
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":3306)/" + dbName

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	rows, err := db.Table(id).Rows()
	if err != nil {
		log.Println("Failed to query table:", err)
		g.IndentedJSON(http.StatusOK, ApiResult{"no_table", err.Error(), nil})
		return
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Println("Failed to get columns:", err)
		return
	}

	// Create a slice to store the results
	var results []map[string]interface{}

	// Iterate over the rows
	for rows.Next() {
		// Create a map to store the current row data
		row := make(map[string]interface{})

		// Create a slice to hold the values
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		// Scan the current row into the values slice
		if err := rows.Scan(values...); err != nil {
			log.Println("Failed to scan row:", err)
			return
		}

		// Iterate over the values and add them to the map
		for i, col := range columns {
			val := *(values[i].(*interface{}))
			if bytesVal, ok := val.([]byte); ok {
				row[col] = string(bytesVal)
			} else {
				row[col] = val
			}
		}

		results = append(results, row)
	}

	g.IndentedJSON(http.StatusOK, ApiResult{"success", id + " table is selected", results})
}
