package main

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
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
	mysqlusername = "root"
	mysqlpassword = "Dev127336"
	mysqlhostname = "127.0.0.1:3306"
	mysqldbname   = "demo"
	pgusername    = "admin"
	pgpassword    = "pg123"
	pghostname    = "127.0.0.1:5432"
	pgdbname      = "demo"
)

// getDSN 返回 MySQL 連線字串
func getMySqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlusername, mysqlpassword, mysqlhostname, mysqldbname)
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

func QueryPostgresSQL(sqlQuery2 string) {
	// 連線到 PostgreSQL
	databaseUrl := "postgres://" + pgusername + ":" + pgpassword + "@" + pghostname + "/" + pgdbname
	// fmt.Println(databaseUrl)

	dbpool, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close() //關掉連線

	// 移除 SQL 結尾的分號（pgx 不接受帶分號的查詢）
	sqlQuery2 = sqlQuery2[:len(sqlQuery2)-1]
	// fmt.Println(sqlQuery2)

	// 執行 SQL 查詢
	var bookName = "D" // 替換成你的查詢條件
	rows, err := dbpool.Query(context.Background(), sqlQuery2, bookName)
	if err != nil {
		log.Fatalf("Query PostgreSQL failed: %v\n", err)
	}
	defer rows.Close()
	// 輸出查詢結果
	fmt.Println("Query PostgreSQL Results:")
	for rows.Next() {
		var name string
		var book_name string
		if err := rows.Scan(&name, &book_name); err != nil {
			log.Fatalf("Failed to scan row: %v\n", err)
		}
		fmt.Printf("Name: %s, Book Name: %s\n", name, book_name)
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
	QueryPostgresSQL(sqlQuery2) // 查詢PostgresSQL DB

}
