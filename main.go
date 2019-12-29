package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"tbox_backend/config"
	_ "tbox_backend/docs"
	"tbox_backend/external"
	"tbox_backend/internal/helpers"
	"tbox_backend/internal/services"
	"tbox_backend/internal/stores"
	"tbox_backend/internal/validator"
	"tbox_backend/routers"
)

// @title TBOX Backend API
// @version 1.0
// @description Swagger API for TBOX Backend.
// @BasePath /api
func main() {
	cfg := config.Load()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	db, err := sql.Open("mysql", cfg.MySQL.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"mysql",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	_ = migration.Up()

	smsService := external.NewSmsService(cfg.SmsService.Url)
	userValidator := validator.NewUserValidator()
	userOtpValidator := validator. NewUserOtpValidator()
	userOtpHelper := helpers.NewUserOtpHelper()
	userHelper := helpers.NewUserHelper("")

	sqlxDb := sqlx.NewDb(db, "mysql")
	userStore := stores.NewUserStore(sqlxDb)
	userOtpStore := stores.NewUserOtpStore(sqlxDb)

	userService := services.NewUserService(
		cfg,
		smsService,
		userValidator,
		userOtpValidator,
		userOtpHelper,
		userHelper,
		userStore,
		userOtpStore,
	)

	phoneNumberLimitConfig := cfg.PhoneNumberRateLimit
	phoneNumberLimiter := helpers.NewPhoneNumberRateLimiters(phoneNumberLimitConfig.Limit, phoneNumberLimitConfig.Burst)
	r := routers.NewRouter(userService, userValidator, phoneNumberLimiter)
	r.IndexRouter(router)
	// setup swagger
	url := ginSwagger.URL(cfg.Swagger.Url)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	_ = router.Run(fmt.Sprintf(":%d", cfg.Port))
}
