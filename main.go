package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
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

// 記錄log用
var logger *zap.Logger

func initLogger() {
	file, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writeSyncer := zapcore.AddSync(file)

	encoderConfig := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(
		// zapcore.NewJSONEncoder(encoderConfig),
		// writeSyncer,
		zapcore.NewConsoleEncoder(encoderConfig),                             // 可能需要 ConsoleEncoder 而非 JSONEncoder
		zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), // 確保寫入檔案 + 終端機
		zap.InfoLevel,
	)

	logger = zap.New(core)
	defer logger.Sync() // 保證在程式結束時同步
}

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
	// 創建 zapgorm2 日誌記錄器
	// zapgormLogger := zapgorm2.New(logger)

	//以下 zapgormLogger.Info 方法不能用
	zapgormLogger := zapgorm2.New(logger).LogMode(gormlogger.Info)
	// 使用 zapgormLogger 記錄 SQL 查詢
	zapgormLogger.Info(context.Background(), "執行 SQL 查詢88", sqlQuery)
	zapgormLogger.Info(context.TODO(), "執行 SQL 查詢99", sqlQuery)

	// 連接 MySQL (GORM)
	dsn := getMySqlDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: zapgormLogger, // 設定 GORM 日誌記錄器
	})
	if err != nil {
		log.Fatalf("連線 MySQL 失敗: %v", err)
		logger.Fatal("連線 MySQL 失敗:", zap.Error(err))
		panic("failed to connect database")
	}
	// fmt.Println("成功連線 MySQL!")
	logger.Info("成功連線 MySQL!")

	// 執行帶參數的 SQL 查詢
	var results []map[string]interface{}
	param := "B" // 搜尋條件
	result := db.Raw(sqlQuery, param).Scan(&results)
	if result.Error != nil {
		log.Fatalf("執行 SQL 查詢失敗: %v", result.Error)
		logger.Fatal("執行 SQL 查詢失敗:", zap.Error(result.Error))
	}
	logger.Info("輸入 SQL 查詢", zap.String("SQL Commend: ", sqlQuery))

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
	initLogger()        // 初始化全域 Logger
	defer logger.Sync() // 保證日誌寫入
	r := QueryMySql(sqlQuery)
	fmt.Println(r) //再一次顯示內容
	logger.Info("執行 MySQL SQL 查詢", zap.String("SQL", sqlQuery))
	// logger.Fatal("sql eroor conentx: ", zap.Error(err))

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
