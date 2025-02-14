package test

import (
	"context"
	"testing"
	"time"

	"github.com/erknas/forum/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreatePost(t *testing.T) {
	store := new(MockStorer)

	ctx := context.Background()

	input := model.PostInput{
		Title:           "Test Post",
		Author:          "Bob",
		Content:         "something",
		CommentsAllowed: true,
	}

	expectedPost := model.CustomPost{
		ID:              1,
		Title:           "Test Post",
		Author:          "Bob",
		Content:         "something",
		CreatedAt:       time.Now(),
		CommentsAllowed: true,
	}

	store.On("CreatePost", ctx, input).Return(expectedPost, nil)

	post, err := store.CreatePost(ctx, input)
	require.NoError(t, err)

	assert.Equal(t, expectedPost.ID, post.ID)
	assert.Equal(t, expectedPost.Title, post.Title)
	assert.Equal(t, expectedPost.Author, post.Author)
	assert.Equal(t, expectedPost.Content, post.Content)
	assert.WithinDuration(t, expectedPost.CreatedAt, post.CreatedAt, time.Second)
	assert.Equal(t, expectedPost.CommentsAllowed, post.CommentsAllowed)

	store.AssertCalled(t, "CreatePost", ctx, input)
}

func TestGetPosts(t *testing.T) {
	store := new(MockStorer)

	ctx := context.Background()

	expectedPosts := []model.CustomPost{
		{ID: 1, Title: "Some Title", Author: "Author 1", Content: "Content", CreatedAt: time.Now(), CommentsAllowed: true, Comments: nil},
		{ID: 2, Title: "Title", Author: "Author 2", Content: "something", CreatedAt: time.Now(), CommentsAllowed: false, Comments: nil},
	}

	store.On("GetPosts", mock.Anything).Return(expectedPosts, nil)

	posts, err := store.GetPosts(ctx)
	require.NoError(t, err)

	require.Equal(t, 2, len(posts))

	assert.Equal(t, expectedPosts[0].ID, posts[0].ID)
	assert.Equal(t, expectedPosts[1].ID, posts[1].ID)
	assert.Equal(t, expectedPosts[0].Title, posts[0].Title)
	assert.Equal(t, expectedPosts[1].Title, posts[1].Title)
	assert.Equal(t, expectedPosts[0].Author, posts[0].Author)
	assert.Equal(t, expectedPosts[1].Author, posts[1].Author)
	assert.Equal(t, expectedPosts[0].Content, posts[0].Content)
	assert.Equal(t, expectedPosts[1].Content, posts[1].Content)
	assert.WithinDuration(t, expectedPosts[0].CreatedAt, posts[0].CreatedAt, time.Second)
	assert.WithinDuration(t, expectedPosts[1].CreatedAt, posts[1].CreatedAt, time.Second)
	assert.True(t, posts[0].CommentsAllowed)
	assert.False(t, posts[1].CommentsAllowed)
	assert.Len(t, posts[0].Comments, 0)
	assert.Len(t, posts[1].Comments, 0)

	store.AssertCalled(t, "GetPosts", mock.Anything)
}

func TestGetPostByID(t *testing.T) {
	store := new(MockStorer)

	ctx := context.Background()

	postID := 1

	expectedPost := model.CustomPost{
		ID:              1,
		Title:           "Test Post",
		Author:          "Test Author",
		Content:         "This is a test post",
		CreatedAt:       time.Now(),
		CommentsAllowed: true,
		Comments:        nil,
	}

	store.On("GetPostByID", ctx, postID).Return(expectedPost, nil)

	post, err := store.GetPostByID(ctx, postID)
	require.NoError(t, err)

	assert.Equal(t, expectedPost.ID, post.ID)
	assert.Equal(t, expectedPost.Title, post.Title)
	assert.Equal(t, expectedPost.Author, post.Author)
	assert.Equal(t, expectedPost.Content, post.Content)
	assert.WithinDuration(t, expectedPost.CreatedAt, post.CreatedAt, time.Second)
	assert.Equal(t, expectedPost.CommentsAllowed, post.CommentsAllowed)
	assert.Len(t, post.Comments, 0)

	store.AssertCalled(t, "GetPostByID", ctx, postID)
}

func TestCreateComment(t *testing.T) {
	store := new(MockStorer)

	ctx := context.Background()

	input := model.CustomCommentInput{
		PostID:   1,
		Author:   "Bob Ross",
		Content:  "test comment",
		ParentID: nil,
	}

	expectedComment := model.CustomComment{
		ID:        1,
		Author:    "Bob Ross",
		Content:   "test comment",
		CreatedAt: time.Now(),
		PostID:    1,
		ParentID:  nil,
	}

	store.On("CreateComment", ctx, input).Return(expectedComment, nil)

	comment, err := store.CreateComment(ctx, input)
	require.NoError(t, err)

	assert.Equal(t, expectedComment.ID, comment.ID)
	assert.Equal(t, expectedComment.Author, comment.Author)
	assert.Equal(t, expectedComment.Content, comment.Content)
	assert.WithinDuration(t, expectedComment.CreatedAt, comment.CreatedAt, time.Second)
	assert.Equal(t, expectedComment.PostID, comment.PostID)
	assert.Equal(t, expectedComment.ParentID, comment.ParentID)

	store.AssertCalled(t, "CreateComment", ctx, input)
}

func TestGetCommentByPost(t *testing.T) {
	store := new(MockStorer)

	ctx := context.Background()

	postID := 1
	limit := 10
	offset := 0

	expectedComments := []model.CustomComment{
		{
			ID:        1,
			Author:    "Test Author 1",
			Content:   "This is a test comment 1",
			CreatedAt: time.Now(),
			PostID:    1,
			ParentID:  nil,
		},
		{
			ID:        2,
			Author:    "Test Author 2",
			Content:   "This is a test comment 2",
			CreatedAt: time.Now(),
			PostID:    1,
			ParentID:  nil,
		},
	}

	store.On("GetCommentsByPost", ctx, postID, limit, offset).Return(expectedComments, nil)

	comments, err := store.GetCommentsByPost(ctx, postID, limit, offset)
	require.NoError(t, err)

	require.Equal(t, len(expectedComments), len(comments))
	assert.Equal(t, expectedComments[0].ID, comments[0].ID)
	assert.Equal(t, expectedComments[1].ID, comments[1].ID)
	assert.Equal(t, expectedComments[0].Author, comments[0].Author)
	assert.Equal(t, expectedComments[1].Author, comments[1].Author)
	assert.Equal(t, expectedComments[0].Content, comments[0].Content)
	assert.Equal(t, expectedComments[1].Content, comments[1].Content)
	assert.WithinDuration(t, expectedComments[0].CreatedAt, comments[0].CreatedAt, time.Second)
	assert.WithinDuration(t, expectedComments[1].CreatedAt, comments[1].CreatedAt, time.Second)
	assert.Equal(t, expectedComments[0].PostID, comments[0].PostID)
	assert.Equal(t, expectedComments[1].PostID, comments[1].PostID)
	assert.Equal(t, expectedComments[0].ParentID, comments[0].ParentID)
	assert.Equal(t, expectedComments[1].ParentID, comments[1].ParentID)

	store.AssertCalled(t, "GetCommentsByPost", ctx, postID, limit, offset)
}

func TestCommentReplies(t *testing.T) {
	store := new(MockStorer)

	ctx := context.Background()

	parentID := 1

	expectedReplies := []model.CustomComment{
		{
			ID:        2,
			Author:    "Author 1",
			Content:   "reply to comment 1",
			CreatedAt: time.Now(),
			PostID:    1,
			ParentID:  &parentID,
		},
		{
			ID:        3,
			Author:    "Author 2",
			Content:   "another reply to comment 1",
			CreatedAt: time.Now(),
			PostID:    1,
			ParentID:  &parentID,
		},
	}

	store.On("GetCommentReplies", ctx, parentID).Return(expectedReplies, nil)

	replies, err := store.GetCommentReplies(ctx, parentID)
	require.NoError(t, err)

	require.Equal(t, len(expectedReplies), len(replies))
	assert.Equal(t, expectedReplies[0].ID, replies[0].ID)
	assert.Equal(t, expectedReplies[0].ID, replies[0].ID)
	assert.Equal(t, expectedReplies[1].Author, replies[1].Author)
	assert.Equal(t, expectedReplies[1].Author, replies[1].Author)
	assert.Equal(t, expectedReplies[1].Content, replies[1].Content)
	assert.Equal(t, expectedReplies[1].Content, replies[1].Content)
	assert.WithinDuration(t, expectedReplies[0].CreatedAt, replies[0].CreatedAt, time.Second)
	assert.WithinDuration(t, expectedReplies[1].CreatedAt, replies[1].CreatedAt, time.Second)
	assert.Equal(t, expectedReplies[0].PostID, replies[0].PostID)
	assert.Equal(t, expectedReplies[1].PostID, replies[1].PostID)
	assert.Equal(t, expectedReplies[0].ParentID, replies[0].ParentID)
	assert.Equal(t, expectedReplies[1].ParentID, replies[1].ParentID)

	store.AssertCalled(t, "GetCommentReplies", ctx, parentID)
}
