package service

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/model"
	"zblog-backend/internal/repository"
	"zblog-backend/internal/utils"
)

type ArticleService interface {
	Create(req *dto.CreateArticleRequest) (*model.Article, error)
	Update(id uint, req *dto.UpdateArticleRequest) (*model.Article, error)
	Delete(id uint) error
	GetByID(id uint) (*model.Article, error)
	GetBySlug(slug string) (*dto.ArticleDetail, error)
	GetList(query *dto.ArticleListQuery) ([]model.Article, int64, error)
	GetAnnouncements() ([]model.Article, error)
}

type articleService struct {
	articleRepo  repository.ArticleRepository
	categoryRepo repository.CategoryRepository
	tagRepo      repository.TagRepository
}

func NewArticleService(
	articleRepo repository.ArticleRepository,
	categoryRepo repository.CategoryRepository,
	tagRepo repository.TagRepository,
) ArticleService {
	return &articleService{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
	}
}

func (s *articleService) Create(req *dto.CreateArticleRequest) (*model.Article, error) {
	html, err := utils.RenderMarkdown(req.ContentMD)
	if err != nil {
		return nil, fmt.Errorf("render markdown: %w", err)
	}

	slug := req.Slug
	if slug == "" {
		slug = utils.GenerateSlug(req.Title)
	}

	baseSlug := slug
	for i := 2; ; i++ {
		existing, err := s.articleRepo.FindBySlug(slug)
		if err != nil {
			break
		}
		_ = existing
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}

	article := &model.Article{
		Title:       req.Title,
		Slug:        slug,
		ContentMD:   req.ContentMD,
		ContentHTML: html,
		Excerpt:     req.Excerpt,
		Cover:       req.Cover,
		Status:      req.Status,
		IsTop:        req.IsTop,
		IsAnnouncement: req.IsAnnouncement,
		CategoryID:   req.CategoryID,
		Series:      req.Series,
		SeriesOrder: req.SeriesOrder,
	}

	if len(req.TagIDs) > 0 {
		tags, err := s.tagRepo.FindByIDs(req.TagIDs)
		if err != nil {
			return nil, fmt.Errorf("find tags: %w", err)
		}
		article.Tags = tags
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"title": article.Title,
		"slug":  article.Slug,
	}).Info("article created")

	return s.articleRepo.FindByID(article.ID)
}

func (s *articleService) Update(id uint, req *dto.UpdateArticleRequest) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	html, err := utils.RenderMarkdown(req.ContentMD)
	if err != nil {
		return nil, fmt.Errorf("render markdown: %w", err)
	}

	article.Title = req.Title
	if req.Slug != "" {
		if existing, err := s.articleRepo.FindBySlug(req.Slug); err == nil && existing.ID != id {
			return nil, fmt.Errorf("slug '%s' 已被其他文章使用", req.Slug)
		}
		article.Slug = req.Slug
	}
	article.ContentMD = req.ContentMD
	article.ContentHTML = html
	article.Excerpt = req.Excerpt
	article.Cover = req.Cover
	article.Status = req.Status
	article.IsTop = req.IsTop
	article.IsAnnouncement = req.IsAnnouncement
	article.CategoryID = req.CategoryID
	article.Series = req.Series
	article.SeriesOrder = req.SeriesOrder

	if req.TagIDs != nil {
		tags, err := s.tagRepo.FindByIDs(req.TagIDs)
		if err != nil {
			return nil, fmt.Errorf("find tags: %w", err)
		}
		if err := s.articleRepo.Update(article); err != nil {
			return nil, err
		}
		if err := model.DB.Model(article).Association("Tags").Replace(tags); err != nil {
			return nil, err
		}
	} else {
		if err := s.articleRepo.Update(article); err != nil {
			return nil, err
		}
	}

	logrus.WithFields(logrus.Fields{
		"title": article.Title,
		"id":    article.ID,
	}).Info("article updated")

	return s.articleRepo.FindByID(article.ID)
}

func (s *articleService) Delete(id uint) error {
	logrus.WithField("id", id).Info("article deleted")
	return s.articleRepo.Delete(id)
}

func (s *articleService) GetByID(id uint) (*model.Article, error) {
	return s.articleRepo.FindByID(id)
}

func (s *articleService) GetBySlug(slug string) (*dto.ArticleDetail, error) {
	article, err := s.articleRepo.FindBySlug(slug)
	if err != nil {
		return nil, err
	}
	if err := s.articleRepo.IncrementViewCount(article.ID); err != nil {
		logrus.WithField("id", article.ID).Warn("increment view count failed")
	}
	article.ViewCount++

	detail := &dto.ArticleDetail{Article: article}

	if article.Series != "" {
		prev, err := s.articleRepo.FindPrevInSeries(article.Series, article.SeriesOrder)
		if err == nil && prev != nil {
			detail.PrevInSeries = &dto.SeriesLink{Slug: prev.Slug, Title: prev.Title}
		}
		next, err := s.articleRepo.FindNextInSeries(article.Series, article.SeriesOrder)
		if err == nil && next != nil {
			detail.NextInSeries = &dto.SeriesLink{Slug: next.Slug, Title: next.Title}
		}

		allArticles, err := s.articleRepo.FindBySeries(article.Series)
		if err == nil {
			detail.SeriesTotal = len(allArticles)
			for _, a := range allArticles {
				detail.SeriesArticles = append(detail.SeriesArticles, dto.SeriesLink{
					Slug:  a.Slug,
					Title: a.Title,
				})
			}
		}
	}

	return detail, nil
}

func (s *articleService) GetList(query *dto.ArticleListQuery) ([]model.Article, int64, error) {
	return s.articleRepo.FindList(query)
}

func (s *articleService) GetAnnouncements() ([]model.Article, error) {
	return s.articleRepo.FindAnnouncements()
}
