package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetAll(context.Context) ([]Post, error)
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64,PaginationFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
		CreateAndInvite(ctx context.Context ,user *User,token string,exp time.Duration) error
	}
	Comments interface {
		GetByPostID(context.Context, int64) ([]Comment, error)
		Create(context.Context, *Comment) error
	}
	Followers interface {
		Follow(ctx context.Context,followerID,userID int64) error
		Unfollow(ctx context.Context,followerID,userID int64) error
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Posts:    &PostStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}

func withTx(db *sql.DB,ctx context.Context,fn func(*sql.Tx) error) error {
	tx,err:=db.BeginTx(ctx, nil)
	if err!=nil{
		return err
	}
	if err:=fn(tx);err!=nil{
		_=tx.Rollback()
		return err
	}
	return tx.Commit()
} 
