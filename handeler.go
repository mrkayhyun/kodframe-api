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

var db *gorm.DB

func InitDB() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	// 환경변수에서 포트 가져오기
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// MySQL 데이터베이스 연결 설정
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":3306)/" + dbName

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

func getApiHandler(g *gin.Context) {
	id := g.Param("id")

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

func getTablesHandler(g *gin.Context) {
	// Show tables query
	rows, err := db.Raw("SHOW TABLES").Rows()
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
	defer rows.Close()

	// Create a map to store the query results
	var tables []Table

	// Iterate over the rows and scan into the map
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		tables = append(tables, Table{TableName: tableName})
	}
	g.IndentedJSON(http.StatusOK, ApiResult{"success", "", tables})
}

func getTableHandler(g *gin.Context) {
	id := g.Param("id")

	rows, err := db.Raw("desc " + id).Rows()
	if err != nil {
		log.Fatalf("테이블 상세정보 조회시 에러: %v", err)
	}

	var columns []TableDesc

	// Iterate over rows and populate TableDesc
	for rows.Next() {
		var desc TableDesc
		if err := rows.Scan(&desc.Field, &desc.Type, &desc.Null, &desc.Key, &desc.Default, &desc.Extra); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		columns = append(columns, desc)
	}

	g.IndentedJSON(http.StatusOK, ApiResult{"success", id + " table desc", columns})
}
