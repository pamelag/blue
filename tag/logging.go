package tag

import (
	"log"
	"time"

	"github.com/pamelag/blue/content"
)

type loggingService struct {
	next Service
}

// NewLoggingService returns a new innstance of a logging service
func NewLoggingService(s Service) Service {
	return &loggingService{s}
}

func (s *loggingService) GetTag(name string, date string) (tag *content.Tag, err error) {
	defer func(begin time.Time) {
		log.Println("method", "GetTag", "took", time.Since(begin), "err", err)
	}(time.Now())
	return s.next.GetTag(name, date)
}
