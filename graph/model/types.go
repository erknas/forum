package model

import (
	"strconv"
	"time"

	"github.com/erknas/forum/pkg/conv"
)

const layout = "02.01.2006 15:04"

type CustomPost struct {
	ID              int              `json:"id"`
	Title           string           `json:"title"`
	Author          string           `json:"author"`
	Content         string           `json:"content"`
	CreatedAt       time.Time        `json:"createdAt"`
	CommentsAllowed bool             `json:"commentsAllowed"`
	Comments        []*CustomComment `json:"comments,omitempty"`
}

type CustomComment struct {
	ID        int       `json:"id"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	PostID    int       `json:"postId"`
	ParentID  *int      `json:"parentId,omitempty"`
}

type CustomCommentInput struct {
	PostID   int    `json:"postId"`
	Author   string `json:"author"`
	Content  string `json:"content"`
	ParentID *int   `json:"parentId,omitempty"`
}

func (p CustomPost) Convert() Post {
	return Post{
		ID:              strconv.Itoa(p.ID),
		Title:           p.Title,
		Author:          p.Author,
		Content:         p.Content,
		CreatedAt:       p.CreatedAt.Format(layout),
		CommentsAllowed: p.CommentsAllowed,
	}
}

func (c CustomComment) Convert() Comment {
	if c.ParentID == nil {
		return Comment{
			ID:        strconv.Itoa(c.ID),
			Author:    c.Author,
			Content:   c.Content,
			CreatedAt: c.CreatedAt.Format(layout),
			PostID:    strconv.Itoa(c.PostID),
			ParentID:  nil,
		}
	}

	parentID := strconv.Itoa(*c.ParentID)

	return Comment{
		ID:        strconv.Itoa(c.ID),
		Author:    c.Author,
		Content:   c.Content,
		CreatedAt: c.CreatedAt.Format(layout),
		PostID:    strconv.Itoa(c.PostID),
		ParentID:  &parentID,
	}
}

func (c CommentInput) Convert() (CustomCommentInput, error) {
	postID, err := conv.ID(c.PostID)
	if err != nil {
		return CustomCommentInput{}, err
	}

	if c.ParentID == nil {
		return CustomCommentInput{
			PostID:   postID,
			Author:   c.Author,
			Content:  c.Content,
			ParentID: nil,
		}, nil
	}

	parentID, err := conv.ID(*c.ParentID)
	if err != nil {
		return CustomCommentInput{}, err
	}

	return CustomCommentInput{
		PostID:   postID,
		Author:   c.Author,
		Content:  c.Content,
		ParentID: &parentID,
	}, nil

}
