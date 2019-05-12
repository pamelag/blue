package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/pamelag/blue/tag"
)

type tagHandler struct {
	t tag.Service
}

func (h *tagHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Get("/{tagName}/{date}", h.getTag)

	return r
}

func (h *tagHandler) getTag(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	name := chi.URLParam(r, "tagName")
	date := chi.URLParam(r, "date")

	tag, err := h.t.GetTag(name, date)
	if err != nil {
		encodeError(ctx, err, w)
		return
	}

	var response = struct {
		Tag         string   `json:"tag"`
		Count       int      `json:"count"`
		Articles    []string `json:"articles"`
		RelatedTags []string `json:"related_tags"`
	}{
		Tag:         tag.Name,
		Count:       tag.ArticleCount,
		Articles:    transformArticleIDs(tag.ArticleIDs),
		RelatedTags: tag.RelatedTags,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		encodeError(ctx, err, w)
		return
	}
}

func transformArticleIDs(ids []int) []string {
	articleIDs := make([]string, 0)
	for _, number := range ids {
		id := strconv.Itoa(number)
		articleIDs = append(articleIDs, id)
	}
	return articleIDs
}
