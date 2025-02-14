package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.64

import (
	"context"

	"github.com/erknas/forum/graph/model"
)

// CreatePost is the resolver for the CreatePost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input model.PostInput) (*model.Post, error) {
	post, err := r.Svc.CreatePost(ctx, input)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// CreateComment is the resolver for the CreateComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.CommentInput) (*model.Comment, error) {
	comment, err := r.Svc.CreateComment(ctx, input)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// GetPosts is the resolver for the GetPosts field.
func (r *queryResolver) GetPosts(ctx context.Context) ([]*model.Post, error) {
	posts, err := r.Svc.Posts(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPostByID is the resolver for the GetPostByID field.
func (r *queryResolver) GetPostByID(ctx context.Context, id string, page *int32, pageSize *int32) (*model.Post, error) {
	post, err := r.Svc.PostByID(ctx, id)
	if err != nil {
		return nil, err
	}

	comments, err := r.Svc.CommentsByPost(ctx, id, page, pageSize)
	if err != nil {
		return nil, err
	}

	post.Comments = comments

	return post, nil
}

// CommentAdded is the resolver for the CommentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *model.Comment, error) {
	ch := r.Sub.Subscribe(postID)

	go func() {
		<-ctx.Done()
		r.Sub.Unsubscribe(postID, ch)
	}()

	return ch, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
