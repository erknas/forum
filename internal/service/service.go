package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/erknas/forum/graph/model"
	"github.com/erknas/forum/internal/storage"
	"github.com/erknas/forum/pkg/conv"
	"github.com/erknas/forum/pkg/pagination"
	"github.com/erknas/forum/pkg/sl"
)

type Servicer interface {
	CreatePost(context.Context, model.PostInput) (*model.Post, error)
	Posts(context.Context) ([]*model.Post, error)
	PostByID(context.Context, string) (*model.Post, error)
	CreateComment(context.Context, model.CommentInput) (*model.Comment, error)
	CommentsByPost(context.Context, string, *int32, *int32) ([]*model.Comment, error)
}

type Service struct {
	store storage.Storer
}

func New(store storage.Storer) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) CreatePost(ctx context.Context, input model.PostInput) (*model.Post, error) {
	if errors := input.ValidatePostInput(); len(errors) > 0 {
		return nil, fmt.Errorf("invalid request data: %v", errors)
	}

	customPost, err := s.store.CreatePost(ctx, input)
	if err != nil {
		slog.Error("failed to create post", sl.Err(err), "input", input)
		return nil, err
	}

	post := customPost.Convert()

	slog.Info("CreatePost OK", "post", post)

	return &post, nil
}

func (s *Service) Posts(ctx context.Context) ([]*model.Post, error) {
	customPosts, err := s.store.GetPosts(ctx)
	if err != nil {
		slog.Error("failed to get posts", sl.Err(err))
		return nil, err
	}

	posts := make([]*model.Post, 0, len(customPosts))

	for _, customPost := range customPosts {
		post := customPost.Convert()
		posts = append(posts, &post)
	}

	slog.Info("Posts OK", "post count", len(posts))

	return posts, nil
}

func (s *Service) PostByID(ctx context.Context, strID string) (*model.Post, error) {
	id, err := conv.ID(strID)
	if err != nil {
		return nil, err
	}

	customPost, err := s.store.GetPostByID(ctx, id)
	if err != nil {
		slog.Error("failed to get post", sl.Err(err), "id", id)
		return nil, err
	}

	post := customPost.Convert()

	slog.Info("PostByID OK", "post", post)

	return &post, nil
}

func (s *Service) CreateComment(ctx context.Context, input model.CommentInput) (*model.Comment, error) {
	if errors := input.ValidateCommentInput(); len(errors) > 0 {
		return nil, fmt.Errorf("invalid request data: %v", errors)
	}

	customInput, err := input.Convert()
	if err != nil {
		slog.Error("failed to convert to custom comment input", sl.Err(err))
		return nil, err
	}

	customComment, err := s.store.CreateComment(ctx, customInput)
	if err != nil {
		slog.Error("failed to create comment", sl.Err(err))
		return nil, err
	}

	comment := customComment.Convert()

	slog.Info("CreateComment OK", "comment", comment)

	return &comment, nil
}

func (s *Service) CommentsByPost(ctx context.Context, id string, page *int32, pageSize *int32) ([]*model.Comment, error) {
	limit, offset := pagination.New(page, pageSize)

	comments, err := s.getComments(ctx, id, offset, limit)
	if err != nil {
		slog.Error("failed to get comments for post", sl.Err(err), "post_id", id)
		return nil, err
	}

	slog.Info("CommentsByPost OK", "post_id", id)

	return comments, nil
}

func (s *Service) getComments(ctx context.Context, strID string, offset int, limit int) ([]*model.Comment, error) {
	id, err := conv.ID(strID)
	if err != nil {
		return nil, err
	}

	customComments, err := s.store.GetCommentsByPost(ctx, id, offset, limit)
	if err != nil {
		return nil, err
	}

	comments := make([]*model.Comment, 0, len(customComments))

	for _, customComment := range customComments {
		comment := customComment.Convert()
		comment.Replies = nil
		comments = append(comments, &comment)
	}

	for _, comment := range comments {
		if err := s.populateReplies(ctx, comment); err != nil {
			return nil, err
		}
	}

	return comments, nil
}

func (s *Service) getCommentReplies(ctx context.Context, strID string) ([]*model.Comment, error) {
	id, err := conv.ID(strID)
	if err != nil {
		return nil, err
	}

	customComments, err := s.store.GetCommentReplies(ctx, id)
	if err != nil {
		slog.Error("failed to get comment replies", sl.Err(err), "comment_id", id)
		return nil, err
	}

	comments := make([]*model.Comment, 0, len(customComments))

	for _, customComment := range customComments {
		comment := customComment.Convert()
		comments = append(comments, &comment)
	}

	return comments, nil
}

func (s *Service) populateReplies(ctx context.Context, comment *model.Comment) error {
	replies, err := s.getCommentReplies(ctx, comment.ID)
	if err != nil {
		return err
	}

	comment.Replies = replies

	for _, reply := range replies {
		if err := s.populateReplies(ctx, reply); err != nil {
			return err
		}
	}

	return nil
}
