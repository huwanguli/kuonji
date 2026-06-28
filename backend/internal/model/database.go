package model

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"zblog-backend/internal/config"
	"zblog-backend/internal/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.DatabaseConfig, logLevel gormlogger.LogLevel) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.NewGormLogger(logLevel),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	if err := db.AutoMigrate(
		&User{},
		&Category{},
		&Tag{},
		&Article{},
		&Comment{},
	); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	DB = db
	logrus.Info("database connected and migrated successfully")
	return nil
}

func SetDB(db *gorm.DB) {
	DB = db
}
