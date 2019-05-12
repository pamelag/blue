package postgres

import (
	"context"
	"log"
	"time"

	"github.com/pamelag/blue/content"
	"github.com/jackc/pgx"
)

const (
	insertArticleStmt = "insert_article"
	insertTagStmt = "insert_tag"
)


type articleRepository struct {
	db *pgx.ConnPool
}

// NewArticleRepository creates a new instance of articleRepository
func NewArticleRepository(pool *pgx.ConnPool) content.ArticleRepository {
	ar := &articleRepository{
		db: pool,
	}

	return ar
}

// Store saves and article and its tags to the articleRepository
func (ar *articleRepository) Store(article *content.Article) error {
	db := ar.db
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second * 2)
	defer cancelFunc()

	tx, err := db.BeginEx(ctx, &pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	articleSQL := "insert into article(title, body, created_on) values($1, $2, $3) returning id"
	articleStmt, err := tx.Prepare(insertArticleStmt, articleSQL)
	if err != nil {
		return err
	}

	err = tx.QueryRow(articleStmt.Name, article.Title, article.Body, article.CreatedOn).Scan(&article.ID)
	if err != nil {
		return err
	}
	
	tagSQL := "insert into tag(article_id, created_on, tag_name, related_tags) values($1, $2, $3, $4)"
	_, err = tx.Prepare(insertTagStmt, tagSQL)
	if err != nil {
		return err
	}

	batch := tx.BeginBatch()

	for _, tag := range article.Tags {
		batch.Queue(insertTagStmt, []interface{}{article.ID, article.CreatedOn, tag.Name, tag.RelatedTags}, nil, nil)
	}

	err = batch.Send(ctx, nil)
	if err != nil {
   if e := tx.Rollback(); e != nil {
      log.Println(e)
   }
 
   // closing batch operation due to error on send
   if e := batch.Close(); e != nil {
		log.Println(e)
   }
   return err
	}
 

	err = batch.Close()
	if err != nil {
		if e := tx.Rollback(); e != nil {
			log.Println(e)
		}
		return err
	}


	return tx.Commit()
}

// GetArticle fetches an Article from the articleRepository
func (ar *articleRepository) GetArticle(id int) (*content.Article, error) {
	db := ar.db
	articleSQL := "select a.id, a.title, a.body, a.created_on, t.tag_name from article a left outer join tag t on a.id = t.article_id where a.id = $1"
	article := &content.Article{}

	rows, err := db.Query(articleSQL, id)
	if err != nil {
		return article, err
	}

	defer rows.Close()
	tags := make([]content.Tag, 0)
	for rows.Next() {
		var tagName *string
		rows.Scan(&article.ID, &article.Title, &article.Body, &article.CreatedOn, &tagName)

		if tagName != nil {
			tag := content.Tag{Name: *tagName}
			tags = append(tags, tag)
		}

	}

	article.Tags = tags
	return article, nil
}
