package test

import (
	"context"

	"github.com/erknas/forum/graph/model"
	"github.com/stretchr/testify/mock"
)

type MockStorer struct {
	mock.Mock
}

func (m *MockStorer) CreatePost(ctx context.Context, input model.PostInput) (model.CustomPost, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(model.CustomPost), args.Error(1)
}

func (m *MockStorer) GetPosts(ctx context.Context) ([]model.CustomPost, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.CustomPost), args.Error(1)
}

func (m *MockStorer) GetPostByID(ctx context.Context, id int) (model.CustomPost, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.CustomPost), args.Error(1)
}

func (m *MockStorer) CreateComment(ctx context.Context, input model.CustomCommentInput) (model.CustomComment, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(model.CustomComment), args.Error(1)
}

func (m *MockStorer) GetCommentsByPost(ctx context.Context, postID, limit, offset int) ([]model.CustomComment, error) {
	args := m.Called(ctx, postID, limit, offset)
	return args.Get(0).([]model.CustomComment), args.Error(1)
}

func (m *MockStorer) GetCommentReplies(ctx context.Context, parentID int) ([]model.CustomComment, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]model.CustomComment), args.Error(1)
}
