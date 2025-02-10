package main

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed **/*.sql
var embedFiles embed.FS

type SqlFiles struct {
	embedFiles embed.FS
}

// MySQL 連線資訊
const (
	mysqlusername = "root"
	mysqlpassword = "Dev127336"
	mysqlhostname = "127.0.0.1:3306"
	mysqldbname   = "demo"
	pgusername    = "admin"
	pgpassword    = "pg123"
	pgdbname      = "demo"
)

// getDSN 返回 MySQL 連線字串
func getMySqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlusername, mysqlpassword, mysqlhostname, mysqldbname)
}

// 通用函數：獲取 PostgreSQL DSN
func getPostgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Taipei",
		"127.0.0.1", pgusername, pgpassword, pgdbname)
}

func QueryMySql(sqlQuery string) []map[string]interface{} {
	// 連接 MySQL (GORM)
	dsn := getMySqlDSN()
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
	return results
}

// QueryPostgresSQL 查詢 PostgreSQL (pgx)
func QueryPostgresSQL(sqlQuery2 string) {
	dsn := getPostgresDSN()
	dbpool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("無法連接到 PostgreSQL: %v", err)
	}
	defer dbpool.Close()

	sqlQuery2 = sqlQuery2[:len(sqlQuery2)-1] // 移除分號
	rows, err := dbpool.Query(context.Background(), sqlQuery2, "D")
	if err != nil {
		log.Fatalf("PostgreSQL 查詢失敗: %v", err)
	}
	defer rows.Close()

	fmt.Println("PostgreSQL (pgx) 查詢結果:")
	for rows.Next() {
		var name string
		var book_name string
		if err := rows.Scan(&name, &book_name); err != nil {
			log.Fatalf("結果掃描失敗: %v", err)
		}
		fmt.Printf("Name: %s, Book Name: %s\n", name, book_name)
	}
}

// QueryPostgresSQL2 查詢 PostgreSQL (GORM)
func QueryPostgresSQL2(sqlQuery2 string) {
	dsn := getPostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("無法連接到 PostgreSQL: %v", err)
	}

	var results []struct {
		Name     string `gorm:"column:name"`
		BookName string `gorm:"column:book_name"`
	}
	if err := db.Raw(sqlQuery2, "D").Scan(&results).Error; err != nil {
		log.Fatalf("PostgreSQL (GORM) 查詢失敗: %v", err)
	}

	fmt.Println("PostgreSQL (GORM) 查詢結果:")
	for _, result := range results {
		fmt.Printf("Name: %s, Book Name: %s\n", result.Name, result.BookName)
	}
}

func main() {
	fmt.Println("Hello, world!")
	// 讀取sql file功能
	s := SqlFiles{embedFiles: embedFiles}
	f, err := s.embedFiles.ReadFile("sql/mysql_demo.sql")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(f))
	sqlQuery := string(f)

	//mysql 功能
	r := QueryMySql(sqlQuery)
	fmt.Println(r) //再一次顯示內容

	s2 := SqlFiles{embedFiles: embedFiles}
	f2, err := s2.embedFiles.ReadFile("sql/pg_demo.sql")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(f2))
	sqlQuery2 := string(f2)
	QueryPostgresSQL(sqlQuery2)  // 查詢PostgresSQL DB
	QueryPostgresSQL2(sqlQuery2) // GORM 方法

}
