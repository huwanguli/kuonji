package repository

import (
	"gorm.io/gorm"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/model"
)

type ArticleRepository interface {
	Create(article *model.Article) error
	Update(article *model.Article) error
	Delete(id uint) error
	FindByID(id uint) (*model.Article, error)
	FindBySlug(slug string) (*model.Article, error)
	FindList(query *dto.ArticleListQuery) ([]model.Article, int64, error)
	FindAnnouncements() ([]model.Article, error)
	FindPrevInSeries(series string, order int) (*model.Article, error)
	FindNextInSeries(series string, order int) (*model.Article, error)
	FindBySeries(series string) ([]model.Article, error)
	RenameSeries(oldName, newName string) error
	RemoveSeries(name string) error
	IncrementViewCount(id uint) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepository) Update(article *model.Article) error {
	return r.db.Save(article).Error
}

func (r *articleRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("article_id = ?", id).Delete(&model.Comment{}).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM article_tags WHERE article_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Article{}, id).Error
	})
}

func (r *articleRepository) FindByID(id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.Preload("Category").Preload("Tags").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) FindBySlug(slug string) (*model.Article, error) {
	var article model.Article
	err := r.db.Preload("Category").Preload("Tags").Where("slug = ?", slug).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) FindList(query *dto.ArticleListQuery) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	db := r.db.Model(&model.Article{}).Preload("Category").Preload("Tags")

	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.CategoryID > 0 {
		db = db.Where("category_id = ?", query.CategoryID)
	}
	if query.TagID > 0 {
		db = db.Where("id IN (SELECT article_id FROM article_tags WHERE tag_id = ?)", query.TagID)
	}
	if query.Keyword != "" {
		db = db.Where("title LIKE ?", "%"+query.Keyword+"%")
	}
	if query.Series != "" {
		db = db.Where("series = ?", query.Series)
	}
	if query.IsTop != nil {
		db = db.Where("is_top = ?", *query.IsTop)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := query.Page
	pageSize := query.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	err := db.Order("is_announcement DESC, is_top DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&articles).Error
	return articles, total, err
}

func (r *articleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *articleRepository) FindAnnouncements() ([]model.Article, error) {
	var articles []model.Article
	err := r.db.Where("is_announcement = ? AND status = ?", 1, 1).
		Order("created_at DESC").
		Find(&articles).Error
	return articles, err
}

func (r *articleRepository) FindPrevInSeries(series string, order int) (*model.Article, error) {
	var article model.Article
	err := r.db.Where("series = ? AND series_order < ? AND status = 1", series, order).
		Order("series_order DESC").
		First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) FindNextInSeries(series string, order int) (*model.Article, error) {
	var article model.Article
	err := r.db.Where("series = ? AND series_order > ? AND status = 1", series, order).
		Order("series_order ASC").
		First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) FindBySeries(series string) ([]model.Article, error) {
	var articles []model.Article
	err := r.db.Where("series = ? AND status = 1", series).
		Order("series_order ASC, created_at ASC").
		Find(&articles).Error
	return articles, err
}

func (r *articleRepository) RenameSeries(oldName, newName string) error {
	return r.db.Model(&model.Article{}).Where("series = ?", oldName).Update("series", newName).Error
}

func (r *articleRepository) RemoveSeries(name string) error {
	return r.db.Model(&model.Article{}).Where("series = ?", name).Update("series", "").Error
}
