package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/erknas/forum/graph/model"
)

type InMemoryStorage struct {
	mu       sync.RWMutex
	posts    map[int]*model.CustomPost
	comments map[int]*model.CustomComment

	postID    int
	commentID int
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		posts:    make(map[int]*model.CustomPost),
		comments: make(map[int]*model.CustomComment),
	}
}

func (s *InMemoryStorage) CreatePost(_ context.Context, input model.PostInput) (post model.CustomPost, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.postID++

	post = model.CustomPost{
		ID:              s.postID,
		Title:           input.Title,
		Author:          input.Author,
		Content:         input.Content,
		CreatedAt:       time.Now(),
		CommentsAllowed: input.CommentsAllowed,
	}

	s.posts[post.ID] = &post

	return post, nil
}

func (s *InMemoryStorage) GetPosts(_ context.Context) ([]model.CustomPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	posts := make([]model.CustomPost, 0, len(s.posts))
	for _, post := range s.posts {
		posts = append(posts, *post)
	}

	return posts, nil
}

func (s *InMemoryStorage) GetPostByID(_ context.Context, id int) (model.CustomPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[id]
	if !ok {
		return model.CustomPost{}, fmt.Errorf("post not found")
	}

	return *post, nil
}

func (s *InMemoryStorage) CreateComment(_ context.Context, input model.CustomCommentInput) (comment model.CustomComment, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post, ok := s.posts[input.PostID]
	if !ok {
		return comment, fmt.Errorf("post not found")
	}

	if !post.CommentsAllowed {
		return comment, fmt.Errorf("comments not allowed")
	}

	if input.ParentID != nil && *input.ParentID > len(post.Comments) {
		return comment, fmt.Errorf("comment does not exist")
	}

	if input.ParentID != nil && post.Comments[*input.ParentID-1] == nil {
		return comment, fmt.Errorf("comment does not exist")
	}

	s.commentID++

	comment = model.CustomComment{
		ID:        s.commentID,
		Author:    input.Author,
		Content:   input.Content,
		CreatedAt: time.Now(),
		PostID:    post.ID,
		ParentID:  input.ParentID,
	}

	s.comments[comment.ID] = &comment

	post.Comments = append(post.Comments, &comment)

	return comment, nil
}

func (s *InMemoryStorage) GetCommentsByPost(_ context.Context, postID int, offset int, limit int) ([]model.CustomComment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[postID]
	if !ok {
		return nil, fmt.Errorf("post not found")
	}

	var comments []model.CustomComment

	for _, comment := range post.Comments {
		if comment.ParentID == nil {
			comments = append(comments, model.CustomComment{
				ID:        comment.ID,
				Author:    comment.Author,
				Content:   comment.Content,
				CreatedAt: comment.CreatedAt,
				PostID:    comment.PostID,
				ParentID:  comment.ParentID,
			})
		}
	}

	start := offset
	end := start + limit

	if start > len(comments) {
		return nil, nil
	}

	if end > len(comments) {
		end = len(comments)
	}

	return comments[start:end], nil
}

func (s *InMemoryStorage) GetCommentReplies(_ context.Context, parentID int) ([]model.CustomComment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var replies []model.CustomComment

	for _, comment := range s.comments {
		if comment.ParentID != nil && *comment.ParentID == parentID {
			replies = append(replies, model.CustomComment{
				ID:        comment.ID,
				Author:    comment.Author,
				Content:   comment.Content,
				CreatedAt: comment.CreatedAt,
				PostID:    comment.PostID,
				ParentID:  &parentID,
			})

		}
	}

	return replies, nil
}
