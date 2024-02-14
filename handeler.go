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

// 환경변수에서 DB 접속 정보를 가져와서 커넥션을 만든다.
func InitDB() error {
	// 환경변수 로드
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

// id 라는 이름의 테이블이 존재하는지 체크 후 데이터를 반환한다.
// Todo : 로우수를 제한 하는 로직 넣어야 됨
func getApiHandler(g *gin.Context) {
	id := g.Param("id")

	rows, err := db.Table(id).Rows()
	if err != nil {
		log.Println("Failed to query table:", err)
		g.IndentedJSON(http.StatusOK, ApiResult{"no_table", err.Error(), ApiBody{}})
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Println("Failed to get columns:", err)
		return
	}

	var results []map[string]interface{}
	var rowCount int

	for rows.Next() {
		row := make(map[string]interface{})

		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		if err := rows.Scan(values...); err != nil {
			log.Println("Failed to scan row:", err)
			return
		}

		for i, col := range columns {
			val := *(values[i].(*interface{}))
			if bytesVal, ok := val.([]byte); ok {
				row[col] = string(bytesVal)
			} else {
				row[col] = val
			}
		}

		results = append(results, row)
		rowCount++
	}

	g.IndentedJSON(http.StatusOK, ApiResult{"success", id + " table is selected", ApiBody{rowCount, results}})
}

// 해당 DB의 테이블 명을 반환한다.
func getTablesHandler(g *gin.Context) {
	rows, err := db.Raw("SHOW TABLES").Rows()
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
	defer rows.Close()

	var tables []Table
	var rowCount int

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		tables = append(tables, Table{TableName: tableName})
		rowCount++
	}
	g.IndentedJSON(http.StatusOK, ApiResult{"success", "", ApiBody{rowCount, tables}})
}

// id 라는 이름의 테이블 구조를 반환한다.
func getTableHandler(g *gin.Context) {
	id := g.Param("id")

	rows, err := db.Raw("desc " + id).Rows()
	if err != nil {
		log.Fatalf("테이블 상세정보 조회시 에러: %v", err)
	}

	var columns []TableDesc
	var rowCount int

	for rows.Next() {
		var desc TableDesc
		if err := rows.Scan(&desc.Field, &desc.Type, &desc.Null, &desc.Key, &desc.Default, &desc.Extra); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		columns = append(columns, desc)
		rowCount++
	}

	g.IndentedJSON(http.StatusOK, ApiResult{"success", id + " table desc", ApiBody{rowCount, columns}})
}
