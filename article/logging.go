package article

import (
	"log"
	"time"

	"github.com/pamelag/blue/content"
)

type loggingService struct {
	next Service
}

// NewLoggingService returns a new instance of logging Service
func NewLoggingService(s Service) Service {
	return &loggingService{s}
}

func (s *loggingService) AddArticle(title, body string, tags []string) (articleID int, err error) {
	defer func(begin time.Time) {
		log.Println("method", "AddArticle", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.AddArticle(title, body, tags)
}

func (s *loggingService) GetArticle(id string) (article *content.Article, err error) {
	defer func(begin time.Time) {
		log.Println("method", "GetArticle", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.GetArticle(id)
}
