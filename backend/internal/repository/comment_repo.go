package repository

import (
	"gorm.io/gorm"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/model"
)

type CommentRepository interface {
	Create(comment *model.Comment) error
	FindByID(id uint) (*model.Comment, error)
	FindList(query *dto.CommentListQuery) ([]model.Comment, int64, error)
	FindAllForAdmin(page, pageSize int, status *int) ([]model.Comment, int64, error)
	UpdateStatus(id uint, status int) error
	Delete(id uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) FindList(query *dto.CommentListQuery) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	db := r.db.Model(&model.Comment{}).Where("article_id = ? AND status = 1", query.ArticleID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Order("created_at ASC").Find(&comments).Error
	return comments, total, err
}

func (r *commentRepository) FindAllForAdmin(page, pageSize int, status *int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	db := r.db.Model(&model.Comment{})
	if status != nil {
		db = db.Where("status = ?", *status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments).Error
	return comments, total, err
}

func (r *commentRepository) UpdateStatus(id uint, status int) error {
	return r.db.Model(&model.Comment{}).Where("id = ?", id).Update("status", status).Error
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}
