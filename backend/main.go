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
	"zblog-backend/internal/utils"
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

	migrateArticleHTML()

	authService := service.NewAuthService(repository.NewUserRepository(model.DB))
	if err := authService.InitAdmin(); err != nil {
		logrus.Warnf("init admin failed: %v", err)
	}

	articleRepo := repository.NewArticleRepository(model.DB)
	categoryRepo := repository.NewCategoryRepository(model.DB)
	tagRepo := repository.NewTagRepository(model.DB)
	commentRepo := repository.NewCommentRepository(model.DB)
	seriesRepo := repository.NewSeriesRepository(model.DB)

	articleService := service.NewArticleService(articleRepo, categoryRepo, tagRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	commentService := service.NewCommentService(commentRepo, articleRepo)
	seriesService := service.NewSeriesService(seriesRepo, articleRepo)

	authH := handler.NewAuthHandler(authService)
	articleH := handler.NewArticleHandler(articleService)
	categoryH := handler.NewCategoryHandler(categoryService)
	tagH := handler.NewTagHandler(tagService)
	commentH := handler.NewCommentHandler(commentService)
	uploadH := handler.NewUploadHandler()
	seriesH := handler.NewSeriesHandler(seriesService)

	r := router.Setup(authH, articleH, categoryH, tagH, commentH, uploadH, seriesH)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logrus.Infof("server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		logrus.Fatalf("server failed: %v", err)
	}
}

func migrateArticleHTML() {
	var articles []model.Article
	if err := model.DB.Select("id", "content_md", "content_html").Find(&articles).Error; err != nil {
		logrus.Warnf("migrate article html: query failed: %v", err)
		return
	}

	updated := 0
	for _, a := range articles {
		html, err := utils.RenderMarkdown(a.ContentMD)
		if err != nil {
			logrus.Warnf("migrate article html: render id=%d failed: %v", a.ID, err)
			continue
		}
		if html == a.ContentHTML {
			continue
		}
		if err := model.DB.Model(&model.Article{}).Where("id = ?", a.ID).Update("content_html", html).Error; err != nil {
			logrus.Warnf("migrate article html: update id=%d failed: %v", a.ID, err)
			continue
		}
		updated++
	}
	if updated > 0 {
		logrus.Infof("migrate article html: re-rendered %d articles", updated)
	}
}
