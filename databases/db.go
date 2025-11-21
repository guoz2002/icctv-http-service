package databases

import (
	"fmt"
	"log"
	"os"
	"sync"

	"icctv-http-service/models"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Init 初始化数据库连接（默认 sqlite，可通过 DB_DRIVER=mysql 切换到 MySQL）
func Init() (*gorm.DB, error) {
	var initErr error

	once.Do(func() {
		driver := os.Getenv("DB_DRIVER")
		if driver == "" {
			driver = "sqlite"
		}

		cfg := &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
		}

		var (
			conn *gorm.DB
			err  error
		)

		switch driver {
		case "mysql":
			host := getenvDefault("DB_HOST", "mysql")
			port := getenvDefault("DB_PORT", "3306")
			user := getenvDefault("DB_USER", "icctv")
			pass := getenvDefault("DB_PASS", "icctv_pass")
			name := getenvDefault("DB_NAME", "icctv")

			dsn := fmt.Sprintf(
				"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				user,
				pass,
				host,
				port,
				name,
			)

			conn, err = gorm.Open(gmysql.Open(dsn), cfg)
		default:
			dsn := getenvDefault("DATABASE_PATH", "icctv.db")
			conn, err = gorm.Open(sqlite.Open(dsn), cfg)
		}

		if err != nil {
			initErr = err
			return
		}

		// 按照依赖顺序迁移表：先建没有外键的表，再建有外键依赖的表
		if err := conn.AutoMigrate(
			&models.Adminer{},
			&models.PublicNetConfig{},
			&models.Building{}, // 必须在 OrangePi 前面，因为 OrangePi 有外键关联
			&models.OrangePi{},
		); err != nil {
			initErr = err
			return
		}

		db = conn
	})

	return db, initErr
}

// MustDB 返回数据库实例（若未初始化会 Panic）
func MustDB() *gorm.DB {
	if db != nil {
		return db
	}

	conn, err := Init()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	return conn
}

func getenvDefault(key, defVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defVal
}
