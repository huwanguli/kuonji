package router

import (
	"github.com/gin-gonic/gin"

	"zblog-backend/internal/config"
	"zblog-backend/internal/handler"
	"zblog-backend/internal/middleware"
)

func Setup(
	authH *handler.AuthHandler,
	articleH *handler.ArticleHandler,
	categoryH *handler.CategoryHandler,
	tagH *handler.TagHandler,
	commentH *handler.CommentHandler,
	uploadH *handler.UploadHandler,
) *gin.Engine {
	if config.Cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Cors(), middleware.RequestLogger(), middleware.Recovery())

	r.Static("/uploads", config.Cfg.Upload.Path)

	api := r.Group("/api")
	{
		api.GET("/articles", articleH.GetPublicList)
		api.GET("/articles/:slug", articleH.GetBySlug)
		api.GET("/series", articleH.GetSeriesList)
		api.GET("/announcements", articleH.GetAnnouncements)
		api.GET("/categories", categoryH.GetAll)
		api.GET("/tags", tagH.GetAll)
		api.GET("/comments", commentH.GetList)
		api.POST("/comments", commentH.Create)
	}

	admin := r.Group("/api/admin")
	{
		admin.POST("/login", authH.Login)

		authGroup := admin.Group("")
		authGroup.Use(middleware.Auth())
		{
			authGroup.GET("/profile", authH.Profile)

			authGroup.GET("/articles", articleH.GetAdminList)
			authGroup.GET("/articles/:id", articleH.GetAdminDetail)
			authGroup.POST("/articles", articleH.Create)
			authGroup.PUT("/articles/:id", articleH.Update)
			authGroup.DELETE("/articles/:id", articleH.Delete)

			authGroup.POST("/categories", categoryH.Create)
			authGroup.PUT("/categories/:id", categoryH.Update)
			authGroup.DELETE("/categories/:id", categoryH.Delete)

			authGroup.POST("/tags", tagH.Create)
			authGroup.PUT("/tags/:id", tagH.Update)
			authGroup.DELETE("/tags/:id", tagH.Delete)

			authGroup.GET("/comments", commentH.GetAdminList)
			authGroup.PUT("/comments/:id", commentH.UpdateStatus)
			authGroup.DELETE("/comments/:id", commentH.Delete)

			authGroup.POST("/upload", uploadH.Upload)
		}
	}

	return r
}
