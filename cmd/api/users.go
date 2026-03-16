package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hannanaarif/Social/internal/store"
)

type userKey string

const userCtx userKey = "user"

// getUserHandler godoc
// @Summary		Get a user
// @Description	Get a user profile by ID
// @Tags			users
// @Produce		json
// @Param			userID	path		int	true	"User ID"
// @Success		200		{object}	store.User
// @Failure		404		{object}	error
// @Failure		500		{object}	error
// @Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	// ID := chi.URLParam(r, "userID")
	// userID, err := strconv.ParseInt(ID, 10, 64)
	// if err != nil {
	// 	app.badRequest(w, r, err)
	// 	return
	// }
	// user, err := app.store.Users.GetByID(r.Context(), userID)
	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, store.ErrNotFound):
	// 		app.notFoundResponse(w, r, err)
	// 	default:
	// 		app.InternalServerError(w, r, err)
	// 	}
	// 	return
	// }

	user, ok := getUserFromContext(r.Context())
	if !ok {
		app.InternalServerError(w, r, errors.New("user not found in context"))
		return
	}
	if err := writeJSON(w, http.StatusOK, user); err != nil {
		app.InternalServerError(w, r, err)
	}
}

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

// followUserHandler godoc
// @Summary		Follow a user
// @Description	Follow a user by ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userID	path		int			true	"User ID"
// @Param			payload	body		FollowUser	true	"Follow payload"
// @Success		200		{object}	nil
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Failure		500		{object}	error
// @Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser, ok := getUserFromContext(r.Context())
	if !ok {
		app.InternalServerError(w, r, errors.New("user not found in context"))
		return
	}

	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := app.store.Followers.Follow(r.Context(), followerUser.ID, payload.UserID); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, nil); err != nil {
		app.InternalServerError(w, r, err)
	}
}

// unfollowUserHandler godoc
// @Summary		Unfollow a user
// @Description	Unfollow a user by ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userID	path		int			true	"User ID"
// @Param			payload	body		FollowUser	true	"Unfollow payload"
// @Success		200		{object}	nil
// @Failure		400		{object}	error
// @Failure		404		{object}	error
// @Failure		500		{object}	error
// @Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unfollowedUser, ok := getUserFromContext(r.Context())
	if !ok {
		app.InternalServerError(w, r, errors.New("user not found in context"))
		return
	}

	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}
	if err := app.store.Followers.Unfollow(r.Context(), payload.UserID, unfollowedUser.ID); err != nil {
		app.InternalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, nil); err != nil {
		app.InternalServerError(w, r, err)
	}
}
// ActivateUser godoc
//
// @Summary      Activates/Register a user
// @Description  Activates/Register a user by invitation token
// @Tags         users
// @Produce      json
// @Param        token   path   string  true  "Invitation token"
// @Success      204     {string} string "User activated"
// @Failure      404     {object} error
// @Failure      500     {object} error
// @Security     ApiKeyAuth
// @Router       /users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token:=chi.URLParam(r,"token");
	if token==""{
		app.badRequest(w,r,errors.New("token is required"))
		return
	}

	err:=app.store.Users.Activate(r.Context(),token)
	if err!=nil{
		switch err {
		case store.ErrNotFound:
			app.badRequest(w, r, err)
		default:
			app.InternalServerError(w, r, err)
		}
		return
	}
	if err:=writeJSON(w,http.StatusOK,nil);err!=nil{
		app.InternalServerError(w,r,err)
	}
	

}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ID := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(ID, 10, 64)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.InternalServerError(w, r, err)
			}
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(ctx context.Context) (*store.User, bool) {
	user, ok := ctx.Value(userCtx).(*store.User)
	return user, ok
}
