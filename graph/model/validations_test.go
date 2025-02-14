package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePostInput(t *testing.T) {
	tests := []struct {
		name           string
		input          PostInput
		expectedErrors map[string]interface{}
	}{
		{
			name: "Valid Input",
			input: PostInput{
				Title:           "Valid Title",
				Author:          "Author Name",
				Content:         "Some content.",
				CommentsAllowed: true,
			},
			expectedErrors: map[string]interface{}{},
		},
		{
			name: "Empty Title",
			input: PostInput{
				Title:           "",
				Author:          "Author Name",
				Content:         "Something.",
				CommentsAllowed: true,
			},
			expectedErrors: map[string]interface{}{
				"title": "title length cannot be zero",
			},
		},
		{
			name: "Title Too Long",
			input: PostInput{
				Title:           "This title is way too long and exceeds the maximum length of one hundred characters, which is not allowed.",
				Author:          "Author Name",
				Content:         "This is some content.",
				CommentsAllowed: true,
			},
			expectedErrors: map[string]interface{}{
				"title": "title length cannot be more than 100 symbols",
			},
		},
		{
			name: "Empty Author",
			input: PostInput{
				Title:           "Valid Title",
				Author:          "",
				Content:         "This is some content.",
				CommentsAllowed: true,
			},
			expectedErrors: map[string]interface{}{
				"author": "author name length cannot be zero",
			},
		},
		{
			name: "Empty Content",
			input: PostInput{
				Title:           "Valid Title",
				Author:          "Author Name",
				Content:         "",
				CommentsAllowed: true,
			},
			expectedErrors: map[string]interface{}{
				"content": "content length cannot be zero",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.input.ValidatePostInput()
			assert.Equal(t, tt.expectedErrors, errors)
		})
	}
}

func TestValidateCommentInput(t *testing.T) {
	tests := []struct {
		name           string
		input          CommentInput
		expectedErrors map[string]interface{}
	}{
		{
			name: "Valid Input",
			input: CommentInput{
				PostID:  "123",
				Author:  "Author",
				Content: "This is a valid comment.",
			},
			expectedErrors: map[string]interface{}{},
		},
		{
			name: "Empty PostID",
			input: CommentInput{
				PostID:  "",
				Author:  "Author",
				Content: "This is a valid comment.",
			},
			expectedErrors: map[string]interface{}{
				"postID": "postID cannot be empty",
			},
		},
		{
			name: "Empty Author",
			input: CommentInput{
				PostID:  "123",
				Author:  "",
				Content: "This is a valid comment.",
			},
			expectedErrors: map[string]interface{}{
				"author": "author name length cannot be zero",
			},
		},
		{
			name: "Empty Content",
			input: CommentInput{
				PostID:  "123",
				Author:  "Author",
				Content: "",
			},
			expectedErrors: map[string]interface{}{
				"content": "content length cannot be zero",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.input.ValidateCommentInput()
			assert.Equal(t, tt.expectedErrors, errors)
		})
	}
}
