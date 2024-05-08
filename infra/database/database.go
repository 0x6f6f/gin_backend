package database

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"gin-boilerplate/config"
)

// 单例模式，保存数据库连接
var (
	DB  *gorm.DB
	err error
)

func GetDB() *gorm.DB {
	return DB
}

// DbConnection create database connection
func DbConnection(defaultDSN, masterDSN, replicaDSN config.DSN) error {
	var db = DB

	logMode := viper.GetBool("DB_LOG_MODE")
	debug := viper.GetBool("DEBUG")

	loglevel := logger.Silent
	if logMode {
		loglevel = logger.Info
	}

	// 连接默认数据库以检查目标数据库是否存在, 如果目标数据库不存在则创建它
	default_db, default_err := gorm.Open(postgres.Open(defaultDSN.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(loglevel),
	})
	if default_err != nil {
		fmt.Println(default_err)
	}
	var count int64
	default_db.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", masterDSN.Dbname).Scan(&count)
	if count == 0 {
		default_db.Exec(fmt.Sprintf("CREATE DATABASE %s;", masterDSN.Dbname))
		log.Printf("已成功创建空数据库并链接")
	}

	// 连接主数据库
	db, err = gorm.Open(postgres.Open(masterDSN.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(loglevel),
	})
	if !debug {
		db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{
				postgres.Open(replicaDSN.DSN),
			},
			Policy: dbresolver.RandomPolicy{},
		}))
	}
	if err != nil {
		log.Fatalf("Db connection error for DSN=[" + masterDSN.DSN + "]")
		return err
	}
	DB = db
	return nil
}

// /*数据库操作*/

// /**
//  * @Description: 初始化数据库
//  * @param masterDSN 主数据库配置
//  * @param replicaDSN 生产环境数据库配置
//  * @return error
//  */
// func InitDB(defaultDSN, masterDSN, replicaDSN config.DSN) error {
// 	// 连接默认数据库以检查目标数据库是否存在
// 	default_db, default_err := ConnectDB(defaultDSN)
// 	if default_err != nil {
// 		fmt.Println(default_err)
// 	}

// 	// 检查目标数据库是否存在
// 	var count int64
// 	default_db.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", dbName).Scan(&count)

// 	// 如果目标数据库不存在，则创建它
// 	if count == 0 {
// 		// 创建目标数据库
// 		default_db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName))
// 		target_db, err := ConnectDB(host, port, user, password, dbName)
// 		if err == nil {
// 			fmt.Println("已成功创建空数据库并链接")
// 		}

// 		// 创建所有数据表
// 		err = createAllTables(target_db)
// 		if err == nil {
// 			fmt.Println("已成功创建所有数据表")
// 		}
// 	}

// 	fmt.Println("数据库已初始化完成")
// 	DB, err = ConnectDB(host, port, user, password, dbName)
// 	return err
// }

// // 清空并重新建立数据库，开发用；调用该函数链接数据库时，每次都会重置数据。
// func ReloadDB(host string, port uint, user, password, dbName string) error {
// 	defaultDbName := "postgres"
// 	// 连接默认数据库以检查目标数据库是否存在
// 	db, err := ConnectDB(host, port, user, password, defaultDbName)
// 	if err != nil {
// 		return err
// 	}
// 	// 检查目标数据库是否存在
// 	var count int64
// 	db.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", dbName).Scan(&count)

// 	// 如果目标数据库存在，则删除它
// 	if count != 0 {
// 		result := db.Exec(fmt.Sprintf("DROP DATABASE %s;", dbName))
// 		if result.Error != nil {
// 			panic(result.Error)
// 		}
// 	}

// 	// 创建目标数据库
// 	db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName))
// 	db, err = ConnectDB(host, port, user, password, dbName)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("已成功创建数据库")

// 	// 创建所有数据表
// 	err = createAllTables(db)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("已成功创建所有数据表")

// 	fmt.Println("数据库已初始化完成")
// 	DB, err = ConnectDB(host, port, user, password, dbName)
// 	return err
// }

// func ConnectDB(DSN string) (*gorm.DB, error) {
// 	// 连接到数据库
// 	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
// 	return db, err
// }

// // CloseDB 关闭数据库连接
// func CloseDB(db *gorm.DB) error {
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return err
// 	}
// 	if err := sqlDB.Close(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // 创建全部表
// func createAllTables(db *gorm.DB) error {
// 	// 使用 AutoMigrate 方法来自动创建表
// 	err := db.AutoMigrate(
// 		&models.Zone{},       // 创建没有外键依赖的表
// 		&models.Department{}, // 依赖Zone
// 		&models.User{},       // 可能依赖其他表，如Department
// 		&models.UserProfile{},
// 		&models.WorkLog{},
// 		&models.Customer{},
// 		&models.Contract{},
// 		&models.SystemLog{},
// 	)
// 	// 检查是否有错误
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
