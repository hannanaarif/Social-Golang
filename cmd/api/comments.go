package main

import (
	"net/http"

	"github.com/hannanaarif/Social/internal/store"
)

type createCommentPayload struct {
	PostID  int64  `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

// createCommentHandler godoc
// @Summary		Create a comment
// @Description	Create a new comment on a post
// @Tags			comments
// @Accept			json
// @Produce		json
// @Param			payload	body		createCommentPayload	true	"Comment payload"
// @Success		201		{object}	store.Comment
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Router			/comments [post]
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
