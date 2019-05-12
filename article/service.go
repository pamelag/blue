package article

import (
	"errors"
	"strconv"

	"github.com/pamelag/blue/content"
)

// ErrInvalidArgument is returned when one or more arguments are invalid
var ErrInvalidArgument = errors.New("invalid id")

// ErrArticleNotFound is returned when no article for a given id is found in the repository
var ErrArticleNotFound = errors.New("article not found")

// Service is an interface for working with articles
type Service interface {
	AddArticle(title, body string, tags []string) (int, error)
	GetArticle(id string) (*content.Article, error)
}

type service struct {
	articles content.ArticleRepository
}

// NewService returns a new service
func NewService(articles content.ArticleRepository) Service {
	return &service{
		articles: articles,
	}
}

func (s *service) AddArticle(title, body string, tags []string) (int, error) {
	article, err := content.NewArticle(title, body, tags)
	if err != nil {
		return 0, err
	}

	err = s.articles.Store(article)
	return article.ID, err

}

func (s *service) GetArticle(id string) (*content.Article, error) {
	number, err := strconv.Atoi(id)
	if err != nil {
		return nil, ErrInvalidArgument
	}

	article, err := s.articles.GetArticle(number)
	if err == nil && article.ID == 0 {
		return article, ErrArticleNotFound
	}

	return article, err

}
