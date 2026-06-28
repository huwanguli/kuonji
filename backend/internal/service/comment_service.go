package service

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"zblog-backend/internal/dto"
	"zblog-backend/internal/model"
	"zblog-backend/internal/repository"
)

type CommentService interface {
	Create(req *dto.CreateCommentRequest, ip string) (*model.Comment, error)
	GetList(query *dto.CommentListQuery) ([]model.Comment, int64, error)
	GetAllForAdmin(page, pageSize int, status *int) ([]model.Comment, int64, error)
	UpdateStatus(id uint, status int) error
	Delete(id uint) error
}

type commentService struct {
	commentRepo repository.CommentRepository
	articleRepo repository.ArticleRepository
}

func NewCommentService(commentRepo repository.CommentRepository, articleRepo repository.ArticleRepository) CommentService {
	return &commentService{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
	}
}

func (s *commentService) Create(req *dto.CreateCommentRequest, ip string) (*model.Comment, error) {
	_, err := s.articleRepo.FindByID(req.ArticleID)
	if err != nil {
		return nil, fmt.Errorf("article not found")
	}

	comment := &model.Comment{
		ArticleID: req.ArticleID,
		ParentID:  req.ParentID,
		Author:    req.Author,
		Email:     req.Email,
		Content:   req.Content,
		Status:    1,
		IP:        ip,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"article_id": req.ArticleID,
		"author":     req.Author,
	}).Info("comment created")
	return comment, nil
}

func (s *commentService) GetList(query *dto.CommentListQuery) ([]model.Comment, int64, error) {
	return s.commentRepo.FindList(query)
}

func (s *commentService) GetAllForAdmin(page, pageSize int, status *int) ([]model.Comment, int64, error) {
	return s.commentRepo.FindAllForAdmin(page, pageSize, status)
}

func (s *commentService) UpdateStatus(id uint, status int) error {
	logrus.WithFields(logrus.Fields{
		"id":     id,
		"status": status,
	}).Info("comment status updated")
	return s.commentRepo.UpdateStatus(id, status)
}

func (s *commentService) Delete(id uint) error {
	logrus.WithField("id", id).Info("comment deleted")
	return s.commentRepo.Delete(id)
}
