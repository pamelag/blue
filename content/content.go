package content

import (
	"errors"
	"strings"
	"time"
)

// Article represents a piece of writing
type Article struct {
	ID        int
	Title     string
	CreatedOn time.Time
	Body      string
	Tags      []Tag
}

// Tag represents a label attached to an Article
type Tag struct {
	Name         string
	ArticleCount int
	ArticleIDs   []int
	RelatedTags  []string
}

// NewArticle creates a new article
func NewArticle(title, body string, tags []string) (*Article, error) {
	if isEmpty(title) {
		return nil, errors.New("Empty Article title")
	}

	if isLengthAboveLimit(title) {
		return nil, errors.New("Article title has a limit of 200 characters")
	}

	if isEmpty(body) {
		return nil, errors.New("Empty Article body")
	}

	return &Article{
		Title:     title,
		Body:      body,
		Tags:      getTags(tags),
		CreatedOn: time.Now(),
	}, nil
}

func getTags(tags []string) []Tag {
	tagHolder := make(map[string][]string)

	for _, v := range tags {
		_, added := tagHolder[v]
		if !added {
			relatedTags := make([]string, 0)
			tagHolder[v] = relatedTags
		}
	}

	uniqueTags := make([]string, 0, len(tagHolder))

	for k := range tagHolder {
		uniqueTags = append(uniqueTags, k)
	}

	tagList := make([]Tag, 0)

	for name, list := range tagHolder {
		tag := Tag{Name: name, RelatedTags: list}
		for _, uniqueTagName := range uniqueTags {
			if name != uniqueTagName {
				tag.RelatedTags = append(tag.RelatedTags, uniqueTagName)
			}
		}
		tagList = append(tagList, tag)
	}

	return tagList
}

func isEmpty(field string) bool {
	return len(strings.Trim(field, " ")) == 0
}

func isLengthAboveLimit(field string) bool {
	return len(field) > 200
}

// ArticleRepository provides access to article store
type ArticleRepository interface {
	Store(article *Article) error
	GetArticle(id int) (*Article, error)
}

// TagRepository provides access to tag store
type TagRepository interface {
	GetTag(name string, date time.Time) (*Tag, error)
}
