package repository

import (
	"gorm.io/gorm"

	"zblog-backend/internal/model"
)

type TagRepository interface {
	FindAll() ([]model.Tag, error)
	FindByIDs(ids []uint) ([]model.Tag, error)
	FindByID(id uint) (*model.Tag, error)
	FindBySlug(slug string) (*model.Tag, error)
	Create(tag *model.Tag) error
	Update(tag *model.Tag) error
	Delete(id uint) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) FindAll() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Order("id DESC").Find(&tags).Error
	return tags, err
}

func (r *tagRepository) FindByIDs(ids []uint) ([]model.Tag, error) {
	var tags []model.Tag
	if len(ids) == 0 {
		return tags, nil
	}
	err := r.db.Where("id IN ?", ids).Find(&tags).Error
	return tags, err
}

func (r *tagRepository) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) FindBySlug(slug string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) Update(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

func (r *tagRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}
