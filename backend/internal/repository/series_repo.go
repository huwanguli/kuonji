package repository

import (
	"gorm.io/gorm"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/model"
)

type SeriesRepository interface {
	Create(series *model.Series) error
	Update(series *model.Series) error
	Delete(id uint) error
	FindAll() ([]model.Series, error)
	FindAllWithCount() ([]dto.SeriesInfo, error)
	FindByName(name string) (*model.Series, error)
	FindByID(id uint) (*model.Series, error)
}

type seriesRepository struct {
	db *gorm.DB
}

func NewSeriesRepository(db *gorm.DB) SeriesRepository {
	return &seriesRepository{db: db}
}

func (r *seriesRepository) Create(series *model.Series) error {
	return r.db.Create(series).Error
}

func (r *seriesRepository) Update(series *model.Series) error {
	return r.db.Save(series).Error
}

func (r *seriesRepository) Delete(id uint) error {
	return r.db.Delete(&model.Series{}, id).Error
}

func (r *seriesRepository) FindAll() ([]model.Series, error) {
	var series []model.Series
	err := r.db.Order("name ASC").Find(&series).Error
	return series, err
}

func (r *seriesRepository) FindAllWithCount() ([]dto.SeriesInfo, error) {
	var results []dto.SeriesInfo
	err := r.db.Raw(`
		SELECT s.id, s.name, s.cover, s.description,
		       COALESCE(COUNT(a.id), 0) AS count
		FROM series s
		LEFT JOIN articles a ON a.series = s.name AND a.status = 1
		GROUP BY s.id, s.name, s.cover, s.description
		ORDER BY s.name ASC
	`).Scan(&results).Error
	return results, err
}

func (r *seriesRepository) FindByName(name string) (*model.Series, error) {
	var s model.Series
	err := r.db.Where("name = ?", name).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *seriesRepository) FindByID(id uint) (*model.Series, error) {
	var s model.Series
	err := r.db.First(&s, id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}
