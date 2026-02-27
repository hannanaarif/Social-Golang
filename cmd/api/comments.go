package main

import (
	"net/http"
	"github.com/hannanaarif/Social/internal/store"
)

type createCommentPayload struct {
	PostID  int64  `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload createCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	comment := &store.Comment{
		PostID:  payload.PostID,
		UserID:  1,
		Content: payload.Content,
	}

	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, comment); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}