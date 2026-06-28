package repository

import (
	"gorm.io/gorm"

	"zblog-backend/internal/model"
)

type CategoryRepository interface {
	FindAll() ([]model.Category, error)
	FindByID(id uint) (*model.Category, error)
	FindBySlug(slug string) (*model.Category, error)
	Create(category *model.Category) error
	Update(category *model.Category) error
	Delete(id uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Order("id DESC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) FindBySlug(slug string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}
