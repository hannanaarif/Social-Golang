package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/google/uuid"
	"github.com/hannanaarif/Social/internal/store"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	store.User			"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	//hash the user password
	//Store the user in the database

	if err:=user.Password.Set(payload.Password);err!=nil{
		app.InternalServerError(w, r, err)
		return
	}

	ctx:=r.Context()

	plainToken:=uuid.New().String()

	//store the token in the database

	hash:=sha256.Sum256([]byte(plainToken))
	hashedToken:=hex.EncodeToString(hash[:])
	

	err:=app.store.Users.CreateAndInvite(ctx,user,hashedToken,app.config.mail.exp)
	if err!=nil{
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequest(w, r, err)
			return
		case store.ErrDuplicateUsername:
			app.badRequest(w, r, err)
			return
		default:
			app.InternalServerError(w, r, err)
			return
		}
	}


	if err := app.jsonResponse(w, http.StatusCreated, "User registered"); err != nil {
		app.InternalServerError(w, r, err)
		return
	}

}
