package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/erknas/forum/graph/model"
	"github.com/erknas/forum/internal/config"
	poolcfg "github.com/erknas/forum/pkg/pool-cfg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ctxTimeout time.Duration = time.Second * 5

type PostgresPool struct {
	pool *pgxpool.Pool
}

func NewPostgresPool(ctx context.Context, cfg *config.Config) (*PostgresPool, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	poolCfg, err := poolcfg.New(cfg)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &PostgresPool{pool: pool}, nil
}

func (p *PostgresPool) CreatePost(ctx context.Context, input model.PostInput) (post model.CustomPost, err error) {
	var (
		postID    int
		createdAt time.Time
	)

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return post, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			fmt.Println(commitErr)
			err = commitErr
		}
	}()

	exists, err := p.postAuthorExists(ctx, input.Author)
	if err != nil {
		return post, err
	}

	insertPost := `INSERT INTO post (title, content, comments_allowed, author_id) 
				   VALUES ($1, $2, $3, $4) 
				   RETURNING id, created_at`

	if !exists {
		authorID, err := p.insertPostAuthor(ctx, input.Author)
		if err != nil {
			return post, err
		}

		if err = tx.QueryRow(ctx, insertPost, input.Title, input.Content, input.CommentsAllowed, authorID).Scan(&postID, &createdAt); err != nil {
			return post, err
		}
	} else {
		authorID, err := p.postAuthorID(ctx, input.Author)
		if err != nil {
			return post, err
		}

		if err = tx.QueryRow(ctx, insertPost, input.Title, input.Content, input.CommentsAllowed, authorID).Scan(&postID, &createdAt); err != nil {
			return post, err
		}
	}

	post = model.CustomPost{
		ID:              postID,
		Title:           input.Title,
		Author:          input.Author,
		Content:         input.Content,
		CreatedAt:       createdAt,
		CommentsAllowed: input.CommentsAllowed,
	}

	return post, nil
}

func (p *PostgresPool) GetPosts(ctx context.Context) ([]model.CustomPost, error) {
	query := `SELECT post.id, post.title, post.content, post.created_at, post.comments_allowed, post_author.author_name FROM post 
			  JOIN post_author ON post.author_id = post_author.id 
			  ORDER BY post.created_at`

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.CustomPost

	for rows.Next() {
		post := model.CustomPost{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.CommentsAllowed, &post.Author); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *PostgresPool) GetPostByID(ctx context.Context, id int) (model.CustomPost, error) {
	query := `SELECT post.id, post.title, post.content, post.created_at, post.comments_allowed, post_author.author_name FROM post 
			  JOIN post_author ON post.author_id = post_author.id 
			  WHERE post.id = $1`

	post := model.CustomPost{}
	if err := p.pool.QueryRow(ctx, query, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.CommentsAllowed, &post.Author); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return post, fmt.Errorf("post not found")
		}
		return post, err
	}

	return post, nil
}

func (p *PostgresPool) CreateComment(ctx context.Context, input model.CustomCommentInput) (comment model.CustomComment, err error) {
	var (
		commentID int
		createdAt time.Time
	)

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return comment, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		commitErr := tx.Commit(ctx)
		if commitErr != nil {
			fmt.Println(commitErr)
			err = commitErr
		}
	}()

	commentExists, err := p.commentExists(ctx, input.PostID, input.ParentID)
	if err != nil {
		return comment, err
	}

	if !commentExists && input.ParentID != nil {
		return comment, fmt.Errorf("comment does not exist")
	}

	allowed, err := p.isAllowed(ctx, input.PostID)
	if err != nil {
		return comment, err
	}

	if !allowed {
		return comment, fmt.Errorf("comments not allowed")
	}

	authorExists, err := p.commentAuthorExists(ctx, input.Author)
	if err != nil {
		return comment, err
	}

	insertComment := `INSERT INTO comment (content, post_id, parent_id, author_id) 
					  VALUES ($1, $2, $3, $4) 
					  RETURNING id, created_at`

	if !authorExists {
		authorID, err := p.insertCommentAuthor(ctx, input.Author)
		if err != nil {
			return comment, err
		}

		if err = tx.QueryRow(ctx, insertComment, input.Content, input.PostID, input.ParentID, authorID).Scan(&commentID, &createdAt); err != nil {
			return comment, err
		}
	} else {
		authorID, err := p.commentAuthorID(ctx, input.Author)
		if err != nil {
			return comment, err
		}

		if err = tx.QueryRow(ctx, insertComment, input.Content, input.PostID, input.ParentID, authorID).Scan(&commentID, &createdAt); err != nil {
			return comment, err
		}
	}

	comment = model.CustomComment{
		ID:        commentID,
		Author:    input.Author,
		Content:   input.Content,
		CreatedAt: createdAt,
		PostID:    input.PostID,
		ParentID:  input.ParentID,
	}

	return comment, nil
}

func (p *PostgresPool) GetCommentsByPost(ctx context.Context, postID int, offset int, limit int) ([]model.CustomComment, error) {
	query := `SELECT comment.id, comment.content, comment.created_at, comment.post_id, comment_author.author_name 
			  FROM comment 
			  JOIN comment_author ON comment.author_id = comment_author.id 
			  WHERE comment.post_id = $1
			  AND parent_id IS NULL
			  ORDER BY comment.created_at
			  LIMIT $2 OFFSET $3`

	rows, err := p.pool.Query(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []model.CustomComment

	for rows.Next() {
		comment := model.CustomComment{}
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.PostID, &comment.Author); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (p *PostgresPool) GetCommentReplies(ctx context.Context, parentID int) ([]model.CustomComment, error) {
	query := `SELECT comment.id, comment.content, comment.created_at, comment.post_id, comment.parent_id, comment_author.author_name 
			  FROM comment 
			  JOIN comment_author ON comment.author_id = comment_author.id 
			  WHERE comment.parent_id = $1
			  ORDER BY comment.created_at`

	rows, err := p.pool.Query(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []model.CustomComment

	for rows.Next() {
		comment := model.CustomComment{}
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.PostID, &comment.ParentID, &comment.Author); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (p *PostgresPool) postAuthorExists(ctx context.Context, author string) (bool, error) {
	var (
		exists bool
		query  = `SELECT EXISTS(SELECT 1 FROM post_author WHERE author_name = $1)`
	)

	if err := p.pool.QueryRow(ctx, query, author).Scan(&exists); err != nil {
		return exists, err
	}

	return exists, nil
}

func (p *PostgresPool) postAuthorID(ctx context.Context, author string) (int, error) {
	var (
		id    int
		query = `SELECT id FROM post_author WHERE author_name = $1`
	)

	if err := p.pool.QueryRow(ctx, query, author).Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func (p *PostgresPool) insertPostAuthor(ctx context.Context, author string) (int, error) {
	var (
		id           int
		insertAuthor = `INSERT INTO post_author (author_name) VALUES ($1) RETURNING id`
	)

	if err := p.pool.QueryRow(ctx, insertAuthor, author).Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func (p *PostgresPool) commentAuthorExists(ctx context.Context, author string) (bool, error) {
	var (
		exists bool
		query  = `SELECT EXISTS(SELECT 1 FROM comment_author WHERE author_name = $1)`
	)

	if err := p.pool.QueryRow(ctx, query, author).Scan(&exists); err != nil {
		return exists, err
	}

	return exists, nil
}

func (p *PostgresPool) commentAuthorID(ctx context.Context, author string) (int, error) {
	var (
		id    int
		query = `SELECT id FROM comment_author WHERE author_name = $1`
	)

	if err := p.pool.QueryRow(ctx, query, author).Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func (p *PostgresPool) insertCommentAuthor(ctx context.Context, author string) (int, error) {
	var (
		id           int
		insertAuthor = `INSERT INTO comment_author (author_name) VALUES ($1) RETURNING id`
	)

	if err := p.pool.QueryRow(ctx, insertAuthor, author).Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func (p *PostgresPool) isAllowed(ctx context.Context, id int) (bool, error) {
	var (
		allowed bool
		query   = `SELECT comments_allowed FROM post WHERE id = $1`
	)

	if err := p.pool.QueryRow(ctx, query, id).Scan(&allowed); err != nil {
		return allowed, err
	}

	return allowed, nil
}

func (p *PostgresPool) commentExists(ctx context.Context, postID int, id *int) (bool, error) {
	var (
		exists bool
		query  = `SELECT EXISTS(SELECT 1 FROM comment WHERE post_id = $1 AND id = $2)`
	)

	if err := p.pool.QueryRow(ctx, query, postID, id).Scan(&exists); err != nil {
		return exists, err
	}

	return exists, nil
}
