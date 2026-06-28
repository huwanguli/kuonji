package testutil

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"zblog-backend/internal/config"
	"zblog-backend/internal/handler"
	"zblog-backend/internal/model"
	"zblog-backend/internal/repository"
	"zblog-backend/internal/router"
	"zblog-backend/internal/service"
)

func SetupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Article{}, &model.Comment{})
	require.NoError(t, err)
	model.SetDB(db)

	config.Cfg = &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
			Expire: 24 * time.Hour,
		},
		Upload: config.UploadConfig{
			Path:        "uploads/",
			MaxSize:     10,
			AllowedExts: []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg"},
		},
		Server: config.ServerConfig{Mode: "debug"},
	}

	authService := service.NewAuthService(repository.NewUserRepository(db))
	require.NoError(t, authService.InitAdmin())

	articleRepo := repository.NewArticleRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	commentRepo := repository.NewCommentRepository(db)

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

	return router.Setup(authH, articleH, categoryH, tagH, commentH, uploadH)
}
