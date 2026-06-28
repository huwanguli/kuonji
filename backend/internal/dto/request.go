package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateArticleRequest struct {
	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug"`
	ContentMD   string `json:"content_md" binding:"required"`
	Excerpt     string `json:"excerpt"`
	Cover       string `json:"cover"`
	Status      int    `json:"status"`
	IsTop       int    `json:"is_top"`
	IsAnnouncement int `json:"is_announcement"`
	CategoryID  *uint  `json:"category_id"`
	TagIDs      []uint `json:"tag_ids"`
	Series      string `json:"series"`
	SeriesOrder int    `json:"series_order"`
}

type UpdateArticleRequest struct {
	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug"`
	ContentMD   string `json:"content_md" binding:"required"`
	Excerpt     string `json:"excerpt"`
	Cover       string `json:"cover"`
	Status      int    `json:"status"`
	IsTop       int    `json:"is_top"`
	IsAnnouncement int `json:"is_announcement"`
	CategoryID  *uint  `json:"category_id"`
	TagIDs      []uint `json:"tag_ids"`
	Series      string `json:"series"`
	SeriesOrder int    `json:"series_order"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

type UpdateTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug"`
}

type CreateCommentRequest struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	ParentID  *uint  `json:"parent_id"`
	Author    string `json:"author" binding:"required"`
	Email     string `json:"email"`
	Content   string `json:"content" binding:"required"`
}

type UpdateCommentStatusRequest struct {
	Status int `json:"status" binding:"oneof=0 1 2"`
}

type ArticleListQuery struct {
	Page       int    `form:"page" binding:"omitempty,min=1"`
	PageSize   int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	CategoryID uint   `form:"category_id"`
	TagID      uint   `form:"tag_id"`
	Keyword    string `form:"keyword"`
	Series     string `form:"series"`
	Status     *int   `form:"status"`
	IsTop      *int   `form:"is_top"`
}

type CommentListQuery struct {
	ArticleID uint `form:"article_id"`
	Page      int  `form:"page" binding:"omitempty,min=1"`
	PageSize  int  `form:"page_size" binding:"omitempty,min=1,max=100"`
}
