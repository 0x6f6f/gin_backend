package database

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"gin-boilerplate/controllers"
	"gin-boilerplate/models"
)

var (
	DB  *gorm.DB
	err error
)

// DbConnection create database connection
func DbConnection(masterDSN, replicaDSN string) error {
	var db = DB

	logMode := viper.GetBool("DB_LOG_MODE")
	debug := viper.GetBool("DEBUG")

	loglevel := logger.Silent
	if logMode {
		loglevel = logger.Info
	}

	db, err = gorm.Open(postgres.Open(masterDSN), &gorm.Config{
		Logger: logger.Default.LogMode(loglevel),
	})
	if !debug {
		db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{
				postgres.Open(replicaDSN),
			},
			Policy: dbresolver.RandomPolicy{},
		}))
	}
	if err != nil {
		log.Fatalf("Db connection error for DSN=["+masterDSN+"]")
		return err
	}
	DB = db
	return nil
}

// GetDB connection
func GetDB() *gorm.DB {
	return DB
}

/*数据库操作*/

/**
 * @Description: 初始化数据库
 * @param host 主机地址
 * @param port 端口
 * @param user 用户名
 * @param password 密码
 * @param dbName 数据库名称
 * @return error
 */
 func InitDB(host string, port uint, user, password, dbName string) error {
	defaultDbName := "postgres"
	// 连接默认数据库
	db, err := ConnectDB(host, port, user, password, defaultDbName)
	if err != nil {
		return err
	}
	// 检查目标数据库是否存在
	var count int64
	db.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", dbName).Scan(&count)
	if count == 0 {
		// 如果目标数据库不存在，则创建它
		db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName))
		db, err = ConnectDB(host, port, user, password, dbName)
		if err != nil {
			return err
		}
		fmt.Println("Database created successfully...")
		create_all_tables(db)
		fmt.Println("Tables created successfully...")
		controllers.LogAction(db, 0, "新建数据库和表")
	} else {
		fmt.Println("Database already exists...")
	}
	return nil
}

/**
 * @Description: 连接到数据库
 * @param host 主机地址
 * @param port 端口
 * @param user 用户名
 * @param password 密码
 * @param dbName 数据库名称
 * @return error
 */
func ConnectDB(host string, port uint, user, password, dbName string) (*gorm.DB, error) {
	// 构建连接字符串
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	// 连接到数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// 检查连接是否有错误
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CloseDB 关闭数据库连接
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	return nil
}

// 创建全部表
func create_all_tables(db *gorm.DB) error {
	// 使用 AutoMigrate 方法来自动创建表
	err := db.AutoMigrate(
		&models.User{}, &models.UserProfile{}, &models.SalesDepartment{}, 
		&models.Zone{}, &models.Saler{}, &models.WorkLog{},
		&models.FinanceSpecialist{}, &models.Accountant{}, &models.Customer{},
		&models.Contract{}, &models.SystemLog{},
	)
	// 检查是否有错误
	if err != nil {
		return err
	}
	return nil
}