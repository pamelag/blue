package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pamelag/blue/article"
	"github.com/pamelag/blue/content"
)

type articleHandler struct {
	a article.Service
}

func (h *articleHandler) router() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.addArticle)
	r.Get("/{id}", h.getArticle)

	return r
}

func (h *articleHandler) addArticle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var request struct {
		Title string
		Body  string
		Tags  []string
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		encodeError(ctx, err, w)
		return
	}

	id, err := h.a.AddArticle(request.Title, request.Body, request.Tags)
	if err != nil {
		encodeError(ctx, err, w)
		return
	}

	var response = struct {
		ID int `json:"id"`
	}{
		ID: id,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		encodeError(ctx, err, w)
		return
	}
}

func (h *articleHandler) getArticle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var request struct {
		ID string
	}

	request.ID = chi.URLParam(r, "id")

	/* id, err := strconv.Atoi(request.ID)
	if err != nil {
		encodeError(ctx, article.ErrInvalidArgument, w)
		return
	} */

	a, err := h.a.GetArticle(request.ID)
	if err != nil {
		encodeError(ctx, err, w)
		return
	}

	/* if a.ID == 0 {
		encodeError(ctx, article.ErrArticleNotFound, w)
		return
	} */

	var response = struct {
		ID    int      `json:"id"`
		Title string   `json:"title"`
		Date  string   `json:"date"`
		Body  string   `json:"body"`
		Tags  []string `json:"tags"`
	}{
		ID:    a.ID,
		Title: a.Title,
		Date:  a.CreatedOn.Format("2006-01-02"),
		Body:  a.Body,
		Tags:  getTagNames(a.Tags),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		encodeError(ctx, err, w)
		return
	}
}

func getTagNames(tags []content.Tag) []string {
	tagNames := make([]string, 0)
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}

	return tagNames
}
