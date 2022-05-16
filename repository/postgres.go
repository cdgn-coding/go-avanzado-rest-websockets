package repository

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	models2 "go-rest-websockets/models"
	"log"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func (p *PostgresUserRepository) GetPaginatedPosts(ctx context.Context, size, page int) ([]models2.Post, error) {
	offset := page * size
	rows, err := p.db.QueryContext(ctx, "SELECT id, post_content, created_at, user_id FROM posts LIMIT $1 OFFSET $2", size, offset)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatalf("error closing rows reader %v", err)
		}
	}()

	var posts = make([]models2.Post, 0, size)
	for rows.Next() {
		var post = models2.Post{}
		err := rows.Scan(&post.Id, &post.PostContent, &post.CreatedAt, &post.UserId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil

}

func (p *PostgresUserRepository) DeletePost(ctx context.Context, post *models2.Post) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1 AND user_id = $2", post.Id, post.UserId)
	return err
}

func (p *PostgresUserRepository) UpdatePost(ctx context.Context, post *models2.Post) error {
	_, err := p.db.ExecContext(ctx, "UPDATE posts SET post_content = $1 WHERE id = $2 AND user_id = $3", post.PostContent, post.Id, post.UserId)
	return err
}

func (p *PostgresUserRepository) GetPostById(ctx context.Context, id string) (*models2.Post, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, post_content, created_at, user_id FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatalf("error closing rows reader %v", err)
		}
	}()

	var post = models2.Post{}
	for rows.Next() {
		err := rows.Scan(&post.Id, &post.PostContent, &post.CreatedAt, &post.UserId)
		if err != nil {
			return nil, err
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p PostgresUserRepository) Close() error {
	return p.db.Close()
}

func NewPostgresUserRepository(databaseUrl string) (*PostgresUserRepository, error) {
	open, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil, err
	}

	repository := &PostgresUserRepository{db: open}

	return repository, nil
}

func (p PostgresUserRepository) InsertUser(ctx context.Context, user *models2.User) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

func (p PostgresUserRepository) GetUserById(ctx context.Context, id string) (*models2.User, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatalf("error closing rows reader %v", err)
		}
	}()

	var user = models2.User{}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Email)
		if err != nil {
			return nil, err
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*models2.User, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatalf("error closing rows reader %v", err)
		}
	}()

	var user = models2.User{}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p PostgresUserRepository) InsertPost(ctx context.Context, post *models2.Post) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3)", post.Id, post.PostContent, post.UserId)
	return err
}
