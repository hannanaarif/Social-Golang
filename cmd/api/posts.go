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

// createPostHandler godoc
// @Summary		Create a post
// @Description	Create a new post
// @Tags			posts
// @Accept			json
// @Produce		json
// @Param			payload	body		createPostPayload	true	"Post payload"
// @Success		201		{object}	store.Post
// @Failure		400		{object}	error
// @Failure		500		{object}	error
// @Router			/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload createPostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		// app.errorJSON(w, err, http.StatusBadRequest)
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
		// app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

}

// getPostHandler godoc
// @Summary		Get a post
// @Description	Get a post by ID
// @Tags			posts
// @Produce		json
// @Param			postID	path		int	true	"Post ID"
// @Success		200		{object}	store.Post
// @Failure		404		{object}	error
// @Failure		500		{object}	error
// @Router			/posts/{postID} [get]
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
		// app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
}

// getAllPostsHandler godoc
// @Summary		Get all posts
// @Description	Get all posts
// @Tags			posts
// @Produce		json
// @Success		200	{array}		store.Post
// @Failure		500	{object}	error
// @Router			/posts [get]
func (app *application) getAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.store.Posts.GetAll(r.Context())
	if err != nil {
		app.InternalServerError(w, r, err)
		// writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusOK, posts); err != nil {
		app.InternalServerError(w, r, err)
		// app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
}

// deletePostHandler godoc
// @Summary		Delete a post
// @Description	Delete a post by ID
// @Tags			posts
// @Produce		json
// @Param			postID	path		int	true	"Post ID"
// @Success		200		{object}	map[string]string
// @Failure		404		{object}	error
// @Failure		500		{object}	error
// @Router			/posts/{postID} [delete]
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

// updatePostHandler godoc
// @Summary		Update a post
// @Description	Update a post by ID
// @Tags			posts
// @Accept			json
// @Produce		json
// @Param			postID	path		int					true	"Post ID"
// @Param			payload	body		updatePostPayload	true	"Post payload"
// @Success		200		{object}	store.Post
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Failure		500		{object}	error
// @Router			/posts/{postID} [patch]
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
				app.notFoundResponse(w, r, err)
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
