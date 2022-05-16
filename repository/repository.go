package repository

import (
	"context"
	models2 "go-rest-websockets/models"
)

type Repository interface {
	InsertUser(ctx context.Context, user *models2.User) error
	GetUserById(ctx context.Context, id string) (*models2.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models2.User, error)
	InsertPost(ctx context.Context, post *models2.Post) error
	GetPostById(ctx context.Context, id string) (*models2.Post, error)
	UpdatePost(ctx context.Context, post *models2.Post) error
	DeletePost(ctx context.Context, post *models2.Post) error
	GetPaginatedPosts(ctx context.Context, size, page int) ([]models2.Post, error)
	Close() error
}
