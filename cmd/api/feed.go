package main

import (
	"net/http"

	"github.com/hannanaarif/Social/internal/store"
)

// getUserFeedHandler godoc
// @Summary		Get user feed
// @Description	Get the feed for the current user
// @Tags			feed
// @Produce		json
// @Param			limit	query		int		false	"Limit"
// @Param			offset	query		int		false	"Offset"
// @Param			sort	query		string	false	"Sort"
// @Param			tags	query		[]string	false	"Tags"
// @Param			search	query		string	false	"Search"
// @Success		200		{array}		store.PostWithMetadata
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	//pagination filtering and sorting
	fq := store.PaginationFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := validate.Struct(fq); err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()
	userID := int64(197)
	feed, err := app.store.Posts.GetUserFeed(ctx, userID, fq)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, feed); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}
