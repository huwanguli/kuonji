package main

import (
	"fmt"

	gormlogger "gorm.io/gorm/logger"
	"github.com/sirupsen/logrus"

	"zblog-backend/internal/config"
	"zblog-backend/internal/handler"
	"zblog-backend/internal/logger"
	"zblog-backend/internal/model"
	"zblog-backend/internal/repository"
	"zblog-backend/internal/router"
	"zblog-backend/internal/service"
)

func main() {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	logger.Init(cfg.Log.Level, cfg.Log.File, cfg.Log.MaxAge, cfg.Log.MaxSize)
	defer logger.Close()

	var gormLogLevel gormlogger.LogLevel
	if cfg.Server.Mode == "debug" {
		gormLogLevel = gormlogger.Info
	} else {
		gormLogLevel = gormlogger.Error
	}

	if err := model.InitDB(&cfg.Database, gormLogLevel); err != nil {
		logrus.Fatalf("failed to initialize database: %v", err)
	}

	authService := service.NewAuthService(repository.NewUserRepository(model.DB))
	if err := authService.InitAdmin(); err != nil {
		logrus.Warnf("init admin failed: %v", err)
	}

	articleRepo := repository.NewArticleRepository(model.DB)
	categoryRepo := repository.NewCategoryRepository(model.DB)
	tagRepo := repository.NewTagRepository(model.DB)
	commentRepo := repository.NewCommentRepository(model.DB)

	articleService := service.NewArticleService(articleRepo, categoryRepo, tagRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	commentService := service.NewCommentService(commentRepo, articleRepo)

	authH := handler.NewAuthHandler(authService)
	articleH := handler.NewArticleHandler(articleService)
	categoryH := handler.NewCategoryHandler(categoryService)
	tagH := handler.NewTagHandler(tagService)
	commentH := handler.NewCommentHandler(commentService)
	uploadH := handler.NewUploadHandler()

	r := router.Setup(authH, articleH, categoryH, tagH, commentH, uploadH)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		logrus.Fatalf("server failed: %v", err)
	}
}
