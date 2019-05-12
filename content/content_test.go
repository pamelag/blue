package content

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArticle(t *testing.T) {
	article, err := NewArticle(
		"latest science shows that potato chips are better for you than sugar",
		"some text, potentially containing simple markup about how potato chips are great",
		[]string{"health", "fitness", "science"})

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 3, len(article.Tags), "they should be equal")

	for _, tag := range article.Tags {
		assert.Equal(t, 2, len(tag.RelatedTags), "they should be equal")
		for _, relatedTag := range tag.RelatedTags {
			assert.NotEqual(t, tag.Name, relatedTag, "they should not be equal")
		}
	}

}
