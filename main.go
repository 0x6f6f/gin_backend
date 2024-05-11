package main

import (
	"fmt"
	"gin-boilerplate/config"
	"gin-boilerplate/infra/database"
	"gin-boilerplate/infra/logger"
	"gin-boilerplate/migrations"
	"gin-boilerplate/repository"
	"gin-boilerplate/routers"
	"time"

	"github.com/robfig/cron"
	"github.com/spf13/viper"
)

var c *cron.Cron

func myTask() {
    // 这里执行定时任务的代码
    fmt.Println("Automatically update customer loan intent")
	repository.AutoUpdateCustomerLoanIntent(database.DB)
	fmt.Println("Migrate customer with 0 loan intent to public sea")
	repository.AutoMigrateCustomerToPublicSea(database.DB)
}

func setupCron() {
    c = cron.New()
    c.AddFunc("@every 1d", myTask) // 这里的"@every 1d"表示每天执行一次
    c.Start()
}

func main() {

	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Shanghai")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
	defaultDSN, masterDSN, replicaDSN := config.DbConfiguration()

	if err := database.DbConnection(defaultDSN, masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}
	//later separate migration
	migrations.Migrate()

	router := routers.SetupRoute()

	setupCron()

	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}
