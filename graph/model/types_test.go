package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToCustomCommentInput(t *testing.T) {
	tests := []struct {
		name  string
		input CommentInput
		want  CustomCommentInput
	}{
		{
			name: "Valid PostID with nil ParentID",
			input: CommentInput{
				PostID:   "1",
				Author:   "author1",
				Content:  "content1",
				ParentID: nil,
			},
			want: CustomCommentInput{
				PostID:   1,
				Author:   "author1",
				Content:  "content1",
				ParentID: nil,
			},
		},
		{
			name: "Valid PostID with empty string ParentID",
			input: CommentInput{
				PostID:   "123",
				Author:   "author2",
				Content:  "content2",
				ParentID: stringPtr(""),
			},
			want: CustomCommentInput{
				PostID:   123,
				Author:   "author2",
				Content:  "content2",
				ParentID: nil,
			},
		},
		{
			name: "Valid PostID with valid ParentID",
			input: CommentInput{
				PostID:   "2",
				Author:   "author3",
				Content:  "content3",
				ParentID: stringPtr("4"),
			},
			want: CustomCommentInput{
				PostID:   2,
				Author:   "author3",
				Content:  "content3",
				ParentID: intPtr(4),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := tt.input.Convert()
			assert.NoError(t, err)
			assert.NotEmpty(t, comment)
			assert.Equal(t, tt.want, comment)
		})
	}
}

func TestConvertToCustomCommentInput_FaileCases(t *testing.T) {
	tests := []struct {
		name      string
		input     CommentInput
		wantError error
	}{
		{
			name: "Invalid PostID",
			input: CommentInput{
				PostID:   "-1",
				Author:   "author1",
				Content:  "content1",
				ParentID: nil,
			},
			wantError: fmt.Errorf("invalid ID -1"),
		},
		{
			name: "Invalid ParentID",
			input: CommentInput{
				PostID:   "123",
				Author:   "author2",
				Content:  "content2",
				ParentID: stringPtr("0"),
			},
			wantError: fmt.Errorf("invalid ID 0"),
		},
		{
			name: "Stirng PostID",
			input: CommentInput{
				PostID:   "string_id",
				Author:   "author3",
				Content:  "content3",
				ParentID: stringPtr("2"),
			},
			wantError: fmt.Errorf("invalid ID string_id"),
		},
		{
			name: "Stirng ParentID",
			input: CommentInput{
				PostID:   "1",
				Author:   "author3",
				Content:  "content3",
				ParentID: stringPtr("parentID"),
			},
			wantError: fmt.Errorf("invalid ID parentID"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := tt.input.Convert()
			assert.Error(t, err)
			assert.Empty(t, comment)
			assert.Equal(t, tt.wantError, err)
		})
	}

}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
