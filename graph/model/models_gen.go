// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Comment struct {
	ID        string     `json:"id"`
	Author    string     `json:"author"`
	Content   string     `json:"content"`
	CreatedAt string     `json:"createdAt"`
	PostID    string     `json:"postID"`
	ParentID  *string    `json:"parentID,omitempty"`
	Replies   []*Comment `json:"replies,omitempty"`
}

type CommentInput struct {
	PostID   string  `json:"postID"`
	Author   string  `json:"author"`
	Content  string  `json:"content"`
	ParentID *string `json:"parentID,omitempty"`
}

type Mutation struct {
}

type Post struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Author          string     `json:"author"`
	Content         string     `json:"content"`
	CreatedAt       string     `json:"createdAt"`
	CommentsAllowed bool       `json:"commentsAllowed"`
	Comments        []*Comment `json:"comments,omitempty"`
}

type PostInput struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	Content         string `json:"content"`
	CommentsAllowed bool   `json:"commentsAllowed"`
}

type Query struct {
}

type Subscription struct {
}
