package service

import (
	"errors"
	"strings"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/model"
	"zblog-backend/internal/repository"
)

type SeriesArticleRepo interface {
	FindBySeries(name string) ([]model.Article, error)
	RenameSeries(oldName, newName string) error
	RemoveSeries(name string) error
}

type SeriesService interface {
	Create(name, cover, description string) (*model.Series, error)
	Update(id uint, name, cover, description string) (*model.Series, error)
	Delete(id uint) error
	GetAllWithCount() ([]dto.SeriesInfo, error)
	GetDetail(name string) (*dto.SeriesDetail, error)
}

type seriesService struct {
	seriesRepo  repository.SeriesRepository
	articleRepo SeriesArticleRepo
}

func NewSeriesService(seriesRepo repository.SeriesRepository, articleRepo SeriesArticleRepo) SeriesService {
	return &seriesService{seriesRepo: seriesRepo, articleRepo: articleRepo}
}

func (s *seriesService) Create(name, cover, description string) (*model.Series, error) {
	existing, err := s.seriesRepo.FindByName(name)
	if err == nil && existing != nil {
		return nil, errors.New("系列名称已存在")
	}

	series := &model.Series{
		Name:        name,
		Cover:       cover,
		Description: description,
	}
	if err := s.seriesRepo.Create(series); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "unique") {
			return nil, errors.New("系列名称已存在")
		}
		return nil, err
	}
	return series, nil
}

func (s *seriesService) Update(id uint, name, cover, description string) (*model.Series, error) {
	series, err := s.seriesRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if name != series.Name {
		existing, err := s.seriesRepo.FindByName(name)
		if err == nil && existing != nil && existing.ID != id {
			return nil, errors.New("系列名称已存在")
		}
	}

	oldName := series.Name
	series.Name = name
	series.Cover = cover
	series.Description = description

	if err := s.seriesRepo.Update(series); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "unique") {
			return nil, errors.New("系列名称已存在")
		}
		return nil, err
	}

	if oldName != name {
		return series, s.articleRepo.RenameSeries(oldName, name)
	}

	return series, nil
}

func (s *seriesService) Delete(id uint) error {
	series, err := s.seriesRepo.FindByID(id)
	if err != nil {
		return err
	}
	if err := s.seriesRepo.Delete(id); err != nil {
		return err
	}
	return s.articleRepo.RemoveSeries(series.Name)
}

func (s *seriesService) GetAllWithCount() ([]dto.SeriesInfo, error) {
	return s.seriesRepo.FindAllWithCount()
}

func (s *seriesService) GetDetail(name string) (*dto.SeriesDetail, error) {
	series, err := s.seriesRepo.FindByName(name)
	if err != nil {
		return nil, err
	}

	articles, err := s.articleRepo.FindBySeries(name)
	if err != nil {
		return nil, err
	}

	var result []dto.SeriesArticle
	for _, a := range articles {
		result = append(result, dto.SeriesArticle{
			ID:          a.ID,
			Title:       a.Title,
			Slug:        a.Slug,
			SeriesOrder: a.SeriesOrder,
		})
	}

	return &dto.SeriesDetail{
		Series: dto.SeriesInfo{
			ID:          series.ID,
			Name:        series.Name,
			Cover:       series.Cover,
			Description: series.Description,
			Count:       int64(len(articles)),
		},
		Articles: result,
	}, nil
}
