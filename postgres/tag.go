package postgres

import (
	"time"

	"github.com/jackc/pgx"
	"github.com/pamelag/blue/content"
)

type tagRepository struct {
	db *pgx.ConnPool
}

// NewTagRepository creates a new instance of tagRepository
func NewTagRepository(pool *pgx.ConnPool) content.TagRepository {
	tr := &tagRepository{
		db: pool,
	}

	return tr
}

func (tr *tagRepository) GetTag(name string, date time.Time) (*content.Tag, error) {
	db := tr.db
	tagSQL := "select t.article_id, t.related_tags from tag t where t.tag_name = $1 and t.created_on = $2"
	tag := &content.Tag{Name: name}

	rows, err := db.Query(tagSQL, name, date.Format("2006-01-02"))
	if err != nil {
		return tag, err
	}

	defer rows.Close()
	articleIDs := make([]int, 0)
	relatedTags := make([]string, 0)

	for rows.Next() {
		var articleID int
		var relatedTagArray []string
		rows.Scan(&articleID, &relatedTagArray)
		articleIDs = append(articleIDs, articleID)
		relatedTags = append(relatedTags, relatedTagArray...)
	}

	tag.ArticleIDs = articleIDs
	tag.RelatedTags = getUniqueTags(relatedTags)
	tag.ArticleCount = len(articleIDs)

	return tag, nil

}

func getUniqueTags(tags []string) []string {
	keys := make(map[string]bool)
	relatedTags := []string{}
	for _, entry := range tags {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			relatedTags = append(relatedTags, entry)
		}
	}
	return relatedTags
}
