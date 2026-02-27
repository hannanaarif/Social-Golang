package main

import (
	"net/http"
	"github.com/hannanaarif/Social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	//pagination filtering and sorting
	fq:= store.PaginationFeedQuery{
		Limit: 20,
		Offset: 0,
		Sort: "desc",
	}

	fq,err:=fq.Parse(r)
	if err!=nil{
		app.badRequest(w,r,err)
		return
	}

	if err:=validate.Struct(fq);err!=nil{
		app.badRequest(w,r,err)
		return
	}

	ctx := r.Context()
	userID := int64(68)
	feed, err := app.store.Posts.GetUserFeed(ctx, userID,fq)
	if err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, feed); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}