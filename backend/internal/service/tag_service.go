package service

import (
	"zblog-backend/internal/model"
	"zblog-backend/internal/repository"
	"zblog-backend/internal/utils"
)

type TagService interface {
	GetAll() ([]model.Tag, error)
	Create(name, slug string) (*model.Tag, error)
	Update(id uint, name, slug string) (*model.Tag, error)
	Delete(id uint) error
}

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}

func (s *tagService) GetAll() ([]model.Tag, error) {
	return s.tagRepo.FindAll()
}

func (s *tagService) Create(name, slug string) (*model.Tag, error) {
	if slug == "" {
		slug = utils.GenerateSlug(name)
	}
	tag := &model.Tag{
		Name: name,
		Slug: slug,
	}
	if err := s.tagRepo.Create(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *tagService) Update(id uint, name, slug string) (*model.Tag, error) {
	tag, err := s.tagRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	tag.Name = name
	if slug != "" {
		tag.Slug = slug
	}
	if err := s.tagRepo.Update(tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *tagService) Delete(id uint) error {
	return s.tagRepo.Delete(id)
}
