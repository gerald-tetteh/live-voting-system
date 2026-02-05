package models

import "context"

type Repository[T any] interface {
	Save(ctx context.Context, entity *T) error
	GetById(ctx context.Context, id string) (*T, error)
	UpdateOne(ctx context.Context, id string, entity *T) error
}