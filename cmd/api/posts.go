package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hannanaarif/Social/internal/store"
)

type postKey string

const postCtx postKey = "post"

type createPostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type updatePostPayload struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload createPostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		// writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		// writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  1,
		Tags:    payload.Tags,
	}

	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.InternalServerError(w, r, err)
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.InternalServerError(w, r, err)
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	// postIDParam := chi.URLParam(r, "postID")
	// postID, err := strconv.ParseInt(postIDParam, 10, 64)
	// if err != nil {
	// 	app.InternalServerError(w, r, err)
	// 	// writeJSONError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// post, err := app.store.Posts.GetByID(r.Context(), postID)
	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		writeJSONError(w, http.StatusNotFound, "post not found")
	// 		return
	// 	}

	// 	app.InternalServerError(w, r, err)
	// 	// writeJSONError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
	// fmt.Println("post", post.ID)

	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.InternalServerError(w, r, err)
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// fmt.Println("comments", comments)

	post.Comments = comments

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.InternalServerError(w, r, err)
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.store.Posts.GetAll(r.Context())
	if err != nil {
		app.InternalServerError(w, r, err)
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusOK, posts); err != nil {
		app.InternalServerError(w, r, err)
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	if err := app.store.Posts.Delete(r.Context(), post.ID); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, map[string]string{
		"message": "post deleted successfully",
	}); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload updatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
}

func (app *application) PostcontextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.InternalServerError(w, r, err)
			return
		}

		post, err := app.store.Posts.GetByID(r.Context(), postID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJSONError(w, http.StatusNotFound, "post not found")
				return
			}

			app.InternalServerError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
