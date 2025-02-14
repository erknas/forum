package mocks

import (
	"context"

	"github.com/erknas/forum/graph/model"
	"github.com/stretchr/testify/mock"
)

type MockServicer struct {
	mock.Mock
}

func (m *MockServicer) CreatePost(ctx context.Context, input model.PostInput) (*model.Post, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockServicer) Posts(ctx context.Context) ([]*model.Post, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (m *MockServicer) PostByID(ctx context.Context, id string) (*model.Post, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockServicer) CreateComment(ctx context.Context, input model.CommentInput) (*model.Comment, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*model.Comment), args.Error(1)
}

func (m *MockServicer) CommentsByPost(ctx context.Context, postID string, limit *int32, offset *int32) ([]*model.Comment, error) {
	args := m.Called(ctx, postID, limit, offset)
	return args.Get(0).([]*model.Comment), args.Error(1)
}
