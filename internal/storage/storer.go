package storage

import (
	"context"

	"github.com/erknas/forum/graph/model"
)

type Storer interface {
	CreatePost(context.Context, model.PostInput) (post model.CustomPost, err error)
	GetPosts(context.Context) ([]model.CustomPost, error)
	GetPostByID(context.Context, int) (model.CustomPost, error)
	CreateComment(context.Context, model.CustomCommentInput) (comment model.CustomComment, err error)
	GetCommentsByPost(context.Context, int, int, int) ([]model.CustomComment, error)
	GetCommentReplies(context.Context, int) ([]model.CustomComment, error)
}
