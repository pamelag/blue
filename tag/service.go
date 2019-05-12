package tag

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/pamelag/blue/content"
)

// ErrInvalidArgument is returned when one or more arguments are invalid
var ErrInvalidArgument = errors.New("invalid tag name or date")

// ErrTagNotFound is returned when no tag for a given name and date is found in the repository
var ErrTagNotFound = errors.New("tag not found")

// Service is an interface for tag related operations
type Service interface {
	GetTag(Tag string, Date string) (*content.Tag, error)
}

type service struct {
	tags content.TagRepository
}

// NewService returns a tag service
func NewService(tags content.TagRepository) Service {
	return &service{
		tags: tags,
	}
}

func (s *service) GetTag(name string, date string) (*content.Tag, error) {

	validName := strings.Trim(name, " ")
	if len(validName) <= 0 {
		return nil, ErrInvalidArgument
	}

	if len(date) != 8 {
		return nil, ErrInvalidArgument
	}

	re := regexp.MustCompile("[0-9]+")
	digitParts := re.FindAllString(date, -1)
	if len(digitParts) != 1 {
		return nil, ErrInvalidArgument
	}

	validDate, err := time.Parse("2006-01-02", date[:4]+"-"+date[4:6]+"-"+date[6:8])
	if err != nil {
		return nil, ErrInvalidArgument
	}

	tag, err := s.tags.GetTag(name, validDate)
	if err == nil && tag.ArticleCount == 0 {
		return nil, ErrTagNotFound
	}
	return tag, err
}
