package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// 初始化数据库
func InitDB() (*sql.DB, error) {
	// 从环境变量中数据库连接串
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		// 正式环境则抛出错误
		return nil, fmt.Errorf("请使用`lade env set DATABASE_URL=libsql://xxx.turso.io --app myapp` 命令指定URL路径")
	}
	authToken := os.Getenv("AUTH_TOKEN")
	if authToken == "" {
		return nil, fmt.Errorf("请使用`lade env set AUTH_TOKEN=xxx --app myapp` 命令指定TOKEN")
	}
	// 生成链接字符串
	dsn := fmt.Sprintf("%s?authToken=%s", dbUrl, authToken)
	// 连接数据库
	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, err
	}
	// 配置数据库链接池
	db.SetMaxOpenConns(25)                 // 最大打开的连接数
	db.SetMaxIdleConns(25)                 // 最大空闲的连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 最大连接生命周期

	// 使用ctx测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// 有错误就释放资源
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		// 测试连接
		return nil, err
	}
	log.Println("[INFO]连接数据库成功...")
	return db, nil
}
