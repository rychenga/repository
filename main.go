package main

import (
	"embed"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//go:embed **/*.sql
var embedFiles embed.FS

type SqlFiles struct {
	embedFiles embed.FS
}

// MySQL 連線資訊
const (
	username = "root"
	password = "Dev127336"
	hostname = "127.0.0.1:3306"
	dbname   = "demo"
)

// getDSN 返回 MySQL 連線字串
func getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, hostname, dbname)
}

func main() {
	fmt.Println("Hello, world!")
	// 讀取sql file功能
	s := SqlFiles{embedFiles: embedFiles}
	f, err := s.embedFiles.ReadFile("sql/demo.sql")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(f))
	sqlQuery := string(f)

	// 連接 MySQL (GORM)
	dsn := getDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("連線 MySQL 失敗: %v", err)
	}
	fmt.Println("成功連線 MySQL!")

	// 執行帶參數的 SQL 查詢
	var results []map[string]interface{}
	param := "B" // 搜尋條件
	result := db.Raw(sqlQuery, param).Scan(&results)
	if result.Error != nil {
		log.Fatalf("執行 SQL 查詢失敗: %v", result.Error)
	}

	// 顯示查詢結果
	fmt.Println("查詢結果:")
	for _, row := range results {
		fmt.Println(row)
	}

}
