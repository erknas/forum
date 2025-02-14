package service

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/erknas/forum/graph/model"
	"github.com/erknas/forum/internal/storage/mocks"
	sub "github.com/erknas/forum/internal/subscription/mocks"
	"github.com/erknas/forum/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const layout = "02.01.2006 15:04"

func TestCreatePost(t *testing.T) {
	tests := []struct {
		name    string
		input   model.PostInput
		wantErr error
	}{
		{
			name: "Valid input post",
			input: model.PostInput{
				Title:           "Title",
				Author:          "Bob",
				Content:         "123 something",
				CommentsAllowed: true,
			},
			wantErr: nil,
		},
		{
			name: "Empty Title",
			input: model.PostInput{
				Title:           "",
				Author:          "Bob",
				Content:         "123 something",
				CommentsAllowed: true,
			},
			wantErr: &gqlerror.Error{Message: "invalid request data", Extensions: map[string]interface{}{"title": "title length cannot be zero"}},
		},
		{
			name: "Empty Author",
			input: model.PostInput{
				Title:           "As123",
				Author:          "",
				Content:         "123 something",
				CommentsAllowed: true,
			},
			wantErr: &gqlerror.Error{Message: "invalid request data", Extensions: map[string]interface{}{"author": "author name length cannot be zero"}},
		},
		{
			name: "Empty Content",
			input: model.PostInput{
				Title:           "As123",
				Author:          "Bob",
				Content:         "",
				CommentsAllowed: true,
			},
			wantErr: &gqlerror.Error{Message: "invalid request data", Extensions: map[string]interface{}{"content": "content length cannot be zero"}},
		},
		{
			name: "Empty all",
			input: model.PostInput{
				Title:           "",
				Author:          "",
				Content:         "",
				CommentsAllowed: true,
			},
			wantErr: &gqlerror.Error{Message: "invalid request data", Extensions: map[string]interface{}{"title": "title length cannot be zero", "author": "author name length cannot be zero", "content": "content length cannot be zero"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storerMock := mocks.NewStorer(t)

			if tt.wantErr == nil {
				storerMock.On("CreatePost", mock.Anything, tt.input).Return(model.CustomPost{}, nil)
			}

			s := &Service{store: storerMock}

			post, err := s.CreatePost(context.Background(), tt.input)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				assert.Empty(t, post)
				storerMock.AssertNotCalled(t, "CreatePost", mock.Anything, tt.input)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, post)
			}
		})
	}
}

func TestPosts(t *testing.T) {
	storerMock := mocks.NewStorer(t)
	s := &Service{store: storerMock}

	customPosts := []model.CustomPost{
		{ID: 1, Title: "Title1", Author: "Author1", Content: "Content1", CreatedAt: time.Now(), CommentsAllowed: true},
		{ID: 2, Title: "Title2", Author: "Author2", Content: "Content2", CreatedAt: time.Now(), CommentsAllowed: false},
	}

	storerMock.On("GetPosts", mock.Anything).Return(customPosts, nil)

	posts, err := s.Posts(context.Background())
	require.NoError(t, err)

	assert.Len(t, posts, 2)

	storerMock.AssertExpectations(t)
}

func TestPostByID(t *testing.T) {
	storerMock := mocks.NewStorer(t)
	s := &Service{store: storerMock}

	customPost := model.CustomPost{
		ID:              1,
		Title:           "Title",
		Author:          "Bob",
		Content:         "Something",
		CreatedAt:       time.Now().UTC(),
		CommentsAllowed: true,
		Comments:        nil,
	}

	storerMock.On("GetPostByID", mock.Anything, customPost.ID).Return(customPost, nil)

	post, err := s.PostByID(context.Background(), "1")
	require.NoError(t, err)

	assert.NotEmpty(t, post)
	assert.Equal(t, post.ID, strconv.Itoa(customPost.ID))
	assert.Equal(t, post.Title, customPost.Title)
	assert.Equal(t, post.Content, customPost.Content)

	assert.Equal(t, post.CreatedAt, customPost.CreatedAt.Format(layout))
	assert.Equal(t, post.CommentsAllowed, customPost.CommentsAllowed)
	assert.Equal(t, len(post.Comments), len(customPost.Comments))

	storerMock.AssertExpectations(t)
}

func TestCreateComment(t *testing.T) {
	tests := []struct {
		name               string
		customCommentInput model.CustomCommentInput
		commentInput       model.CommentInput
		expectedComment    model.CustomComment
		comment            model.Comment
		wantErr            error
	}{
		{
			name: "Valid input",
			customCommentInput: model.CustomCommentInput{
				PostID:   1,
				Author:   "Bob",
				Content:  "Something",
				ParentID: nil,
			},
			commentInput: model.CommentInput{
				Author:   "Bob",
				Content:  "Something",
				PostID:   "1",
				ParentID: nil,
			},
			expectedComment: model.CustomComment{
				ID:        1,
				Author:    "Bob",
				Content:   "Something",
				CreatedAt: time.Now(),
				PostID:    1,
				ParentID:  nil,
			},
			comment: model.Comment{
				ID:        "1",
				Author:    "Bob",
				Content:   "Something",
				CreatedAt: time.Now().Format(layout),
				PostID:    "1",
				ParentID:  nil,
			},
			wantErr: nil,
		},
		{
			name: "Empty postID",
			commentInput: model.CommentInput{
				Author:   "Bob",
				Content:  "Something",
				PostID:   "",
				ParentID: nil,
			},
			wantErr: &gqlerror.Error{Message: "invalid request data", Extensions: map[string]interface{}{"postID": "postID cannot be empty"}},
		},
		{
			name: "Empty author",
			commentInput: model.CommentInput{
				Author:   "",
				Content:  "Something",
				PostID:   "2",
				ParentID: nil,
			},
			wantErr: &gqlerror.Error{Message: "invalid request data", Extensions: map[string]interface{}{"author": "author name length cannot be zero"}},
		},
		{
			name: "Empty content",
			commentInput: model.CommentInput{
				Author:   "Bob",
				Content:  "",
				PostID:   "2",
				ParentID: nil,
			},
			wantErr: &gqlerror.Error{Message: "invalid request data", Extensions: map[string]interface{}{"content": "content length cannot be zero"}},
		},
		{
			name: "Empty all",
			commentInput: model.CommentInput{
				Author:   "",
				Content:  "",
				PostID:   "",
				ParentID: nil,
			},
			wantErr: &gqlerror.Error{Message: "invalid request data", Extensions: map[string]interface{}{"postID": "postID cannot be empty", "author": "author name length cannot be zero", "content": "content length cannot be zero"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storerMock := mocks.NewStorer(t)
			subscribeMock := sub.NewSubscriber(t)

			if tt.wantErr == nil {
				storerMock.On("CreateComment", mock.Anything, tt.customCommentInput).Return(tt.expectedComment, nil)
				subscribeMock.On("Publish", tt.comment.PostID, &tt.comment).Return(nil)
			}

			s := &Service{store: storerMock, sub: subscribeMock}

			comment, err := s.CreateComment(context.Background(), tt.commentInput)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
				assert.Empty(t, comment)
				storerMock.AssertNotCalled(t, "CreateComment", mock.Anything, tt.customCommentInput)
				subscribeMock.AssertNotCalled(t, "Publish", tt.commentInput.PostID, &tt.comment)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, comment)
			}
		})
	}

}

func TestCommetsByPost(t *testing.T) {
	storerMock := mocks.NewStorer(t)
	s := &Service{store: storerMock}

	customCommetns := []model.CustomComment{
		{ID: 1, Author: "Author1", Content: "Content1", CreatedAt: time.Now(), PostID: 1},
	}

	var page, pageSie int32 = 1, 10

	limit, offset := pagination.New(&page, &pageSie)

	storerMock.On("GetCommentsByPost", mock.Anything, 1, offset, limit).Return(customCommetns, nil)
	storerMock.On("GetCommentReplies", mock.Anything, 1).Return(nil, nil)

	comments, err := s.CommentsByPost(context.Background(), "1", &page, &pageSie)
	require.NoError(t, err)
	assert.NotEmpty(t, comments)

	storerMock.AssertExpectations(t)
}
