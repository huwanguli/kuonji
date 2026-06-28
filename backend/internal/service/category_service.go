package service

import (
	"zblog-backend/internal/model"
	"zblog-backend/internal/repository"
	"zblog-backend/internal/utils"
)

type CategoryService interface {
	GetAll() ([]model.Category, error)
	Create(name, slug, description string) (*model.Category, error)
	Update(id uint, name, slug, description string) (*model.Category, error)
	Delete(id uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) GetAll() ([]model.Category, error) {
	return s.categoryRepo.FindAll()
}

func (s *categoryService) Create(name, slug, description string) (*model.Category, error) {
	if slug == "" {
		slug = utils.GenerateSlug(name)
	}
	cat := &model.Category{
		Name:        name,
		Slug:        slug,
		Description: description,
	}
	if err := s.categoryRepo.Create(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) Update(id uint, name, slug, description string) (*model.Category, error) {
	cat, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	cat.Name = name
	if slug != "" {
		cat.Slug = slug
	}
	cat.Description = description
	if err := s.categoryRepo.Update(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) Delete(id uint) error {
	return s.categoryRepo.Delete(id)
}
